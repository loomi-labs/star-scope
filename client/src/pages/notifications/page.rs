use js_sys::Date;
use log::debug;
use prost_types::Timestamp;
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;
use wasm_bindgen::JsValue;

use crate::components::messages::create_error_msg_from_status;
use crate::Services;
use crate::services::grpc::{Event};

fn displayTimestamp(option: Option<Timestamp>) -> String {
    if let Some(timestamp) = option {
        let datetime = Date::new(&JsValue::from_f64(timestamp.seconds as f64 * 1000.0));
        let asString = datetime.to_locale_string("en-US", &JsValue::from_str("date"));
        return format!("{}", asString);
    }
    return "".to_string();
}

#[component(inline_props)]
pub fn EventComponent<G: Html>(cx: Scope, event: Event) -> View<G> {
    view! {cx,
        div(class="flex flex-col my-4 p-4 bg-gray-100 dark:bg-gray-800 rounded-lg shadow") {
            div(class="flex flex-row justify-between") {
                div(class="flex flex-col") {
                    p(class="text-lg font-bold") { (event.title.clone()) }
                    p(class="text-sm") { (displayTimestamp(event.timestamp.clone())) }
                    p(class="text-sm") { (event.description.clone()) }
                }
            }
        }
    }
}

#[component]
pub fn Events<G: Html>(cx: Scope) -> View<G> {
    let overview_state = use_context::<OverviewState>(cx);
    let pastEvents = create_memo(cx, || {
        overview_state
            .pastEvents
            .get()
            .iter()
            .take(10)
            .cloned()
            .collect::<Vec<_>>()
    });
    let newEvents = create_memo(cx, || {
        overview_state
            .newEvents
            .get()
            .iter()
            .take(10)
            .cloned()
            .collect::<Vec<_>>()
    });

    view! {cx,
        p(class="text-2xl font-bold") { "Past events" }
        div(class="flex flex-col") {
            Keyed(
                iterable=pastEvents,
                key=|event| event.id.clone(),
                view=|cx,event| {
                    view!{cx,
                        EventComponent(event=event)
                    }
                }
            )
        }
        p(class="text-2xl font-bold") { "New events" }
        div(class="flex flex-col") {
            Keyed(
                iterable=newEvents,
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
    pastEvents: RcSignal<Vec<Event>>,
    newEvents: RcSignal<Vec<Event>>,
}

impl OverviewState {
    pub fn new() -> Self {
        Self {
            pastEvents: create_rc_signal(vec![]),
            newEvents: create_rc_signal(vec![]),
        }
    }
}

async fn query_events(cx: Scope<'_>) {
    let overview_state = use_context::<OverviewState>(cx);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request({});
    let response = services
        .grpc_client
        .get_event_service()
        .list_events(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(response) = response {
        *overview_state.pastEvents.modify() = response.events;
        debug!("Events: {:?}", *overview_state.pastEvents.get());
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
            let mut events = overview_state.newEvents.modify();
            events.push(event);
            *overview_state.newEvents.modify() = events.clone();
        }
    });
}

#[component]
pub async fn Notifications<G: Html>(cx: Scope<'_>) -> View<G> {
    provide_context(cx, OverviewState::new());

    // query_channels(cx.to_owned()).await;
    query_events(cx.to_owned()).await;
    subscribe_to_events(cx.to_owned());

    view! {cx,
        div(class="flex flex-col h-full w-full p-8") {
            h1(class="text-4xl font-bold pb-4") { "Overview" }
            div(class="flex flex-col p-8 bg-white dark:bg-purple-500 rounded-lg shadow") {
                Events {}
            }
        }
    }
}
