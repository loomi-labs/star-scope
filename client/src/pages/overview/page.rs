use log::debug;
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;

use crate::components::messages::create_error_msg_from_status;
use crate::Services;
use crate::services::grpc::{Channel, Event};

#[component(inline_props)]
pub fn EventComponent<G: Html>(cx: Scope, event: Event) -> View<G> {
    view! {cx,
        div(class="flex flex-col my-4 p-4 bg-gray-100 dark:bg-gray-800 rounded-lg shadow") {
            div(class="flex flex-row justify-between") {
                div(class="flex flex-col") {
                    p(class="text-lg font-bold") { (event.title.clone()) }
                    p(class="text-sm") { (event.description.clone()) }
                }
            }
        }
    }
}

#[component]
pub fn Channels<G: Html>(cx: Scope) -> View<G> {
    let overview_state = use_context::<OverviewState>(cx);
    let channels = create_memo(cx, || {
        overview_state
            .channels
            .get()
            .iter()
            .take(10)
            .cloned()
            .collect::<Vec<_>>()
    });
    let events = create_memo(cx, || {
        overview_state
            .events
            .get()
            .iter()
            .take(10)
            .cloned()
            .collect::<Vec<_>>()
    });

    view! {cx,
        // div(class="flex flex-col") {
        //     Keyed(
        //         iterable=channels,
        //         key=|channel| channel.name.clone(),
        //         view=|cx,channel| {
        //             view! {cx,
        //                 div(class="flex flex-col my-4 p-4 bg-gray-100 dark:bg-gray-800 rounded-lg shadow") {
        //                     div(class="flex flex-row justify-between") {
        //                         div(class="flex flex-col") {
        //                             p(class="text-lg font-bold") { (channel.name.clone()) }
        //                             p(class="text-sm") { ("channel description") }
        //                         }
        //                         div(class="flex flex-col justify-center") {
        //                             p(class="text-sm") { "Last message" }
        //                             p(class="text-sm") { ("channel last message") }
        //                         }
        //                     }
        //                     div(class="flex flex-row justify-between") {
        //                         div(class="flex flex-col") {
        //                             p(class="text-sm") { "Nr. of messages" }
        //                             p(class="text-sm") { ("45") }
        //                         }
        //                     }
        //                 }
        //             }
        //         }
        //     )
        // }
        div(class="flex flex-col") {
            Keyed(
                iterable=events,
                key=|event| event.id.clone(),
                view=|cx,event| {
                    view!{cx,
                        EventComponent(event=event)
                    }
                }
            )
        }
    }
}

#[derive(Debug, Clone)]
pub struct OverviewState {
    channels: RcSignal<Vec<Channel>>,
    events: RcSignal<Vec<Event>>,
}

impl OverviewState {
    pub fn new() -> Self {
        Self {
            channels: create_rc_signal(vec![]),
            events: create_rc_signal(vec![]),
        }
    }
}

async fn query_channels(cx: Scope<'_>) {
    let overview_state = use_context::<OverviewState>(cx);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request({});
    let response = services
        .grpc_client
        .get_user_service()
        .list_channels(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(response) = response {
        *overview_state.channels.modify() = response.channels;
        debug!("Channels: {:?}", *overview_state.channels.get());
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

fn subscribe_to_events(cx: Scope) {
    spawn_local_scoped(cx, async move {
        let overview_state = use_context::<OverviewState>(cx);
        let services = use_context::<Services>(cx);
        let mut event_stream = services
            .grpc_client
            .get_event_service()
            .event_stream(services.grpc_client.create_request({}))
            .await
            .unwrap()
            .into_inner();
        while let Some(event) = event_stream.message().await.unwrap() {
            debug!("Received event: {:?}", event);
            let mut events = overview_state.events.modify();
            events.push(event);
            *overview_state.events.modify() = events.clone();
        }
    });
}

#[component]
pub async fn Overview<G: Html>(cx: Scope<'_>) -> View<G> {
    provide_context(cx, OverviewState::new());

    query_channels(cx.to_owned()).await;
    subscribe_to_events(cx.to_owned());

    view! {cx,
        div(class="flex flex-col h-full w-full p-8") {
            h1(class="text-2xl font-bold") { "Overview" }
            div(class="flex flex-col p-8 bg-white dark:bg-gray-600 rounded-lg shadow") {
                Channels {}
            }
        }
    }
}
