use std::str::FromStr;

use chrono::{Duration, NaiveDateTime};
use js_sys::Date;
use log::debug;
use prost_types::Timestamp;
use sycamore::prelude::*;
use wasm_bindgen::JsCast;
use wasm_bindgen::JsValue;
use web_sys::{Event, HtmlSelectElement};
use enum_iterator::{all, Sequence};

use crate::{EventsState, Services};
use crate::components::messages::create_error_msg_from_status;
use crate::services::grpc;

fn display_timestamp(option: Option<Timestamp>) -> String {
    if let Some(timestamp) = option {
        let datetime = Date::new(&JsValue::from_f64(timestamp.seconds as f64 * 1000.0));
        let asString = datetime.to_locale_string("en-US", &JsValue::from_str("date"));
        return format!("{}", asString);
    }
    return "".to_string();
}

// fn get_type_icon(event_type: grpc::EventType) -> String {
//     match event_type {
//         grpc::EventType::Funding => "icon-[ep--coin]".to_string(),
//         grpc::EventType::Staking => "icon-[arcticons--coinstats]".to_string(),
//         grpc::EventType::Dex => "icon-[fluent--money-24-regular]".to_string(),
//         grpc::EventType::Governance => "icon-[icon-park-outline--palace]".to_string(),
//     }.to_string()
// }


#[component(inline_props)]
pub fn EventComponent<G: Html>(cx: Scope, event: grpc::Event) -> View<G> {
    let event_type = event.event_type();

    view! {cx,
        div(class="flex flex-col my-4 p-4 bg-gray-100 dark:bg-purple-700 rounded-lg shadow") {
            div(class="flex flex-row justify-between") {
                div(class="flex flex-row items-center") {
                    div(class="rounded-full h-14 w-14 aspect-square mr-4 bg-gray-300 dark:bg-purple-600 flex items-center justify-center") {
                        img(src=event.chain_image_url, alt="Event Logo", class="h-12 w-12")
                    }
                    div(class="flex flex-col") {
                        p(class="text-lg font-bold") { (event.title.clone()) }
                        p(class="text-sm") { (display_timestamp(event.timestamp.clone())) }
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
    let notifications_state = use_context::<NotificationsState>(cx);
    let events = create_memo(cx, || {
        events_state
            .events
            .get()
            .iter()
            .filter(|_event| {
                let read_status_filter = notifications_state.read_status_filter.get();
                match read_status_filter.as_ref() {
                    ReadStatusFilter::ALL => true,
                    ReadStatusFilter::READ => true,
                    ReadStatusFilter::UNREAD => true,
                }
            }).filter(|event| {
            let time_filter = notifications_state.time_filter.get();
            match time_filter.as_ref().as_time_range() {
                None => true,
                Some((start, end)) => {
                    event.timestamp.clone().unwrap().seconds > start.timestamp() && event.timestamp.clone().unwrap().seconds <= end.timestamp()
                }
            }
        })
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
    read_status_filter: RcSignal<ReadStatusFilter>,
    time_filter: RcSignal<TimeFilter>,
}

impl NotificationsState {
    pub fn new() -> Self {
        Self {
            read_status_filter: create_rc_signal(ReadStatusFilter::default()),
            time_filter: create_rc_signal(TimeFilter::default()),
        }
    }
}

async fn query_events(cx: Scope<'_>) {
    let events_state = use_context::<EventsState>(cx);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(grpc::ListEventsRequest { start_time: None });
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
pub enum EventTypeFilter {
    ALL,
    FUNDING,
    STAKING,
    DEXES,
    GOVERNANCE,
}

#[derive(Debug, Clone, Sequence)]
pub enum ReadStatusFilter {
    ALL,
    READ,
    UNREAD,
}

#[derive(Debug, PartialEq, Eq)]
pub struct ReadStatusFilterError;

impl ReadStatusFilter {
    fn get_filter_from_hash() -> Self {
        let hash = web_sys::window().unwrap().location().hash().unwrap();

        match hash.as_str() {
            "#/read" => ReadStatusFilter::READ,
            "#/unread" => ReadStatusFilter::UNREAD,
            _ => ReadStatusFilter::ALL,
        }
    }

    fn default() -> Self {
        ReadStatusFilter::get_filter_from_hash()
    }

    fn to_string(&self) -> String {
        match self {
            ReadStatusFilter::ALL => "all".to_string(),
            ReadStatusFilter::READ => "read".to_string(),
            ReadStatusFilter::UNREAD => "unread".to_string(),
        }
    }
}

impl FromStr for ReadStatusFilter {
    type Err = ReadStatusFilterError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "all" => Ok(ReadStatusFilter::ALL),
            "read" => Ok(ReadStatusFilter::READ),
            "unread" => Ok(ReadStatusFilter::UNREAD),
            _ => Err(ReadStatusFilterError),
        }
    }
}

const DROPDOWN_DIV_CLASS: &str = "relative inline-flex items-center";
const DROPDOWN_ICON_CLASS: &str = "absolute left-0 top-0 h-full flex items-center pl-2 pointer-events-none text-gray-500 dark:text-purple-600";
const DROPDOWN_SELECT_CLASS: &str = "block capitalize pl-8 pr-4 py-2 rounded border-0 duration-300 hover:bg-sky-400 dark:text-purple-600 dark:bg-purple-700 dark:hover:bg-purple-800";


#[component]
pub fn ReadStatusFilterDropdown<G: Html>(cx: Scope) -> View<G> {
    let notifications_state = use_context::<NotificationsState>(cx);
    let input_ref = create_node_ref(cx);

    let handle_change = |event: Event| {
        let target: HtmlSelectElement = event.target().unwrap().unchecked_into();
        notifications_state
            .read_status_filter
            .set(ReadStatusFilter::from_str(&target.value()).unwrap());
    };

    let options = View::new_fragment(
        all::<ReadStatusFilter>().map(|f| {
            let cloned_f = f.clone();
            view! { cx, option(value=cloned_f.to_string(), class="capitalize") { (f.to_string()) } }
        }).collect()
    );

    view! { cx,
        div(class=DROPDOWN_DIV_CLASS) {
            div(class=DROPDOWN_ICON_CLASS) {
                span(class="icon-[mdi--envelope-outline]")
            }
            select(ref=input_ref,
                class=DROPDOWN_SELECT_CLASS,
                on:change=handle_change,
            ) {
                (options)
            }
        }
    }
}

#[derive(Debug, Clone, Sequence)]
pub enum TimeFilter {
    All,
    Today,
    Yesterday,
    OneWeek,
    OneMonth,
    OneYear,
}

#[derive(Debug, PartialEq, Eq)]
pub struct TimeFilterError;

impl TimeFilter {
    // fn get_filter_from_hash() -> Self {
    //     let hash = web_sys::window().unwrap().location().hash().unwrap();
    //
    //     match hash.as_str() {
    //         "#/read" => TimeFilter::READ,
    //         "#/unread" => TimeFilter::UNREAD,
    //         _ => TimeFilter::ALL,
    //     }
    // }

    fn default() -> Self {
        // TimeFilter::get_filter_from_hash()
        TimeFilter::All
    }

    fn to_string(&self) -> String {
        match self {
            TimeFilter::All => "all".to_string(),
            TimeFilter::Today => "today".to_string(),
            TimeFilter::Yesterday => "yesterday".to_string(),
            TimeFilter::OneWeek => "one week".to_string(),
            TimeFilter::OneMonth => "one month".to_string(),
            TimeFilter::OneYear => "one year".to_string(),
        }
    }

    fn as_time_range(&self) -> Option<(NaiveDateTime, NaiveDateTime)> {
        let js_date = Date::new_0();
        let milliseconds = js_date.get_time();
        let seconds = (milliseconds / 1000.0) as i64;
        let today = NaiveDateTime::from_timestamp_opt(seconds, 0).unwrap().date().and_hms_opt(0, 0, 0).unwrap();
        match self {
            TimeFilter::All => None,
            TimeFilter::Today => Some((today, today + Duration::days(1))),
            TimeFilter::Yesterday => Some((today - Duration::days(1), today)),
            TimeFilter::OneWeek => Some((today - Duration::days(7), today + Duration::days(1))),
            TimeFilter::OneMonth => Some((today - Duration::days(7), today + Duration::days(30))),
            TimeFilter::OneYear => Some((today - Duration::days(7), today + Duration::days(365))),
        }
    }
}

impl FromStr for TimeFilter {
    type Err = TimeFilterError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "all" => Ok(TimeFilter::All),
            "today" => Ok(TimeFilter::Today),
            "yesterday" => Ok(TimeFilter::Yesterday),
            "this_week" => Ok(TimeFilter::OneWeek),
            "this_month" => Ok(TimeFilter::OneMonth),
            "this_year" => Ok(TimeFilter::OneYear),
            _ => Err(TimeFilterError),
        }
    }
}

#[component]
pub fn TimeFilterDropdown<G: Html>(cx: Scope) -> View<G> {
    let notifications_state = use_context::<NotificationsState>(cx);
    let input_ref = create_node_ref(cx);

    let handle_change = |event: Event| {
        let target: HtmlSelectElement = event.target().unwrap().unchecked_into();
        notifications_state
            .time_filter
            .set(TimeFilter::from_str(&target.value()).unwrap());
    };

    let options = View::new_fragment(
        all::<TimeFilter>().map(|f| {
            let cloned_f = f.clone();
            view! { cx, option(value=cloned_f.to_string(), class="capitalize") { (f.to_string()) } }
        }).collect()
    );

    view! { cx,
        div(class=DROPDOWN_DIV_CLASS) {
            div(class=DROPDOWN_ICON_CLASS) {
                span(class="icon-[system-uicons--cubes]")
            }
            select(ref=input_ref,
                class=DROPDOWN_SELECT_CLASS,
                on:change=handle_change,
            ) {
                (options)
            }
        }
    }
}


#[component(inline_props)]
pub async fn Notifications<G: Html>(cx: Scope<'_>, filter: EventTypeFilter) -> View<G> {
    provide_context(cx, NotificationsState::new());

    query_events(cx.to_owned()).await;

    debug!("filter: {:?}", filter);

    view! {cx,
        div(class="flex flex-col") {
            div(class="hidden lg:flex flex-row justify-between items-center") {
                h1(class="text-4xl font-bold pb-4") { "Notifications" }
                div(class="flex flex-row space-x-4 h-8") {
                    ReadStatusFilterDropdown {}
                    TimeFilterDropdown {}
                }
            }
            div(class="lg:hidden flex flex-col") {
                h1(class="text-4xl font-bold pb-4") { "Notifications" }
                div(class="flex flex-wrap") {
                    div(class="w-full sm:w-auto flex-shrink-0 flex-grow-0 mb-4 sm:mb-0 sm:mr-4") {
                        ReadStatusFilterDropdown {}
                    }
                    div(class="w-full sm:w-auto flex-shrink-0 flex-grow-0") {
                        TimeFilterDropdown {}
                    }
                }
            }
            div(class="flex flex-col") {
                Events {}
            }
        }
    }
}
