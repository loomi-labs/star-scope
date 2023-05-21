use js_sys::Date;
use log::debug;
use prost_types::Timestamp;
use sycamore::prelude::*;
use wasm_bindgen::JsValue;

use crate::components::messages::create_error_msg_from_status;
use crate::{EventsState, Services};
use crate::services::grpc::{Event, ListEventsRequest, EventType};

fn displayTimestamp(option: Option<Timestamp>) -> String {
    if let Some(timestamp) = option {
        let datetime = Date::new(&JsValue::from_f64(timestamp.seconds as f64 * 1000.0));
        let asString = datetime.to_locale_string("en-US", &JsValue::from_str("date"));
        return format!("{}", asString);
    }
    return "".to_string();
}

fn getTypeIcon(event_type: EventType) -> String {
    match event_type {
        EventType::Funding => "icon-[ep--coin]".to_string(),
        EventType::Staking => "icon-[arcticons--coinstats]".to_string(),
        EventType::Dex => "icon-[fluent--money-24-regular]".to_string(),
        EventType::Governance => "icon-[icon-park-outline--palace]".to_string(),
    }.to_string()

}


#[component(inline_props)]
pub fn EventComponent<G: Html>(cx: Scope, event: Event) -> View<G> {
    let event_type = event.event_type();

    view! {cx,
        div(class="flex flex-col my-4 p-4 bg-gray-100 dark:bg-purple-700 rounded-lg shadow") {
            div(class="flex flex-row justify-between") {
                div(class="flex flex-row items-center") {
                    div(class="rounded-full h-14 w-14 bg-gray-300 dark:bg-purple-600 flex items-center justify-center mr-4") {
                        img(src=event.chain_image_url, alt="Event Logo", class="h-12 w-12")
                    }
                    div(class="flex flex-col") {
                        p(class="text-lg font-bold") { (event.title.clone()) }
                        p(class="text-sm") { (displayTimestamp(event.timestamp.clone())) }
                        p(class="text-sm") { (event.description.clone()) }
                    }
                }
            }
        }
    }
}

#[component]
pub fn Events<G: Html>(cx: Scope) -> View<G> {
    let events_state = use_context::<EventsState>(cx);
    let events = create_memo(cx, || {
        events_state
            .events
            .get()
            .iter()
            .take(10)
            .cloned()
            .collect::<Vec<_>>()
    });

    view! {cx,
        div(class="flex flex-col") {
            Keyed(
                iterable=events,
                key=|event| event.timestamp.clone(),
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
pub struct NotificationsState {
}

impl NotificationsState {
    pub fn new() -> Self {
        Self {
        }
    }
}

async fn query_events(cx: Scope<'_>) {
    let events_state = use_context::<EventsState>(cx);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(ListEventsRequest{start_time: None});
    let response = services
        .grpc_client
        .get_event_service()
        .list_events(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(response) = response {
        events_state.addEvents(response.events);
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

#[derive(Debug)]
pub enum Filter {
    ALL,
    FUNDING,
    STAKING,
    DEXES,
    GOVERNANCE,
}

#[component(inline_props)]
pub async fn Notifications<G: Html>(cx: Scope<'_>, filter: Filter) -> View<G> {
    provide_context(cx, NotificationsState::new());

    query_events(cx.to_owned()).await;

    debug!("filter: {:?}", filter);

    view! {cx,
        div(class="flex flex-col h-full w-full p-8") {
            h1(class="text-4xl font-bold pb-4") { "Notifications" }
            div(class="flex flex-col") {
                Events {}
            }
        }
    }
}
