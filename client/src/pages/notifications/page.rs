use std::fmt::{Display, Formatter};
use std::str::FromStr;

use chrono::{Duration, NaiveDateTime};
use enum_iterator::{all, Sequence};
use inflector::Inflector;
use js_sys::Date;
use log::debug;
use prost_types::Timestamp;
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;
use wasm_bindgen::closure::Closure;
use wasm_bindgen::JsCast;
use wasm_bindgen::JsValue;
use web_sys::{
    Event, HtmlDivElement, HtmlSelectElement, IntersectionObserver, IntersectionObserverEntry,
};

use crate::components::messages::create_error_msg_from_status;
use crate::types::protobuf::event::EventType;
use crate::types::protobuf::grpc;
use crate::utils::url::{add_or_update_query_params, get_query_param};
use crate::{EventsState, Services};

fn display_timestamp(option: Option<Timestamp>, locale: String) -> String {
    if let Some(timestamp) = option {
        let datetime = Date::new(&JsValue::from_f64(timestamp.seconds as f64 * 1000.0));
        let asString = datetime.to_locale_string(locale.as_ref(), &JsValue::from_str("date"));
        return format!("{}", asString);
    }
    "".to_string()
}

#[component(inline_props)]
pub fn EventBadge<G: Html>(cx: Scope, event_type: EventType) -> View<G> {
    let text_color = match event_type {
        EventType::Funding => "text-badge-green",
        EventType::Staking => "text-badge-red",
        EventType::Dex => "text-badge-orange",
        EventType::Governance => "text-badge-blue",
    };
    let border_color = match event_type {
        EventType::Funding => "border-badge-green",
        EventType::Staking => "border-badge-red",
        EventType::Dex => "border-badge-orange",
        EventType::Governance => "border-badge-blue",
    };
    view! {cx,
        span(class={format!("text-xs rounded-full border flex items-center px-3 mx-8 {} {}", border_color, text_color)}) { (event_type.as_str_name().to_title_case()) }
    }
}

async fn mark_event_as_read(cx: Scope<'_>, event_id: String) {
    let events_state = use_context::<EventsState>(cx);
    let services = use_context::<Services>(cx);
    let request = services
        .grpc_client
        .create_request(grpc::MarkEventReadRequest {
            event_id: event_id.clone(),
        });
    let response = services
        .grpc_client
        .get_event_service()
        .mark_event_read(request)
        .await
        .map(|res| res.into_inner());
    if response.is_ok() {
        events_state.mark_as_read(event_id);
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

#[component(inline_props)]
pub fn EventComponent<G: Html>(cx: Scope, rc_event: RcSignal<grpc::Event>) -> View<G> {
    let event = rc_event.get().as_ref().clone();
    let event_type = event.event_type();

    let is_collapsed = create_signal(cx, true);
    let notifications_state = use_context::<NotificationsState>(cx);
    let locale = notifications_state.locale.get();

    // TODO: make this proper
    let is_clamping = event.description.len() > 250;

    let in_viewport = create_signal(cx, false);
    let event_ref = create_node_ref(cx);

    let boxed = Box::new(
        move |entries: Vec<JsValue>, _observer: IntersectionObserver| {
            for entry in entries {
                let entry: IntersectionObserverEntry = entry.unchecked_into();
                if entry.is_intersecting() {
                    in_viewport.set(true);
                }
            }
        },
    ) as Box<dyn FnMut(Vec<JsValue>, IntersectionObserver)>;
    let handler: Box<dyn FnMut(Vec<JsValue>, IntersectionObserver) + 'static> =
        unsafe { std::mem::transmute(boxed) };
    let callback = Closure::wrap(handler);

    let observer = IntersectionObserver::new(callback.as_ref().unchecked_ref())
        .expect("Failed to create IntersectionObserver");

    callback.forget(); // Prevent the closure from being dropped prematurely

    on_mount(cx, move || {
        if let Ok(element) = event_ref
            .get::<DomNode>()
            .unchecked_into::<HtmlDivElement>()
            .dyn_into::<web_sys::Element>()
        {
            observer.observe(&element);
        }
    });

    on_cleanup(cx, move || {
        if let Ok(_element) = event_ref
            .get::<DomNode>()
            .unchecked_into::<HtmlDivElement>()
            .dyn_into::<web_sys::Element>()
        {
            // TODO: ("Call observer.unobserve(element) here (or observer.disconnect()");
        }
    });

    let rc_event_cloned = rc_event.clone();
    create_effect(cx, move || {
        if *in_viewport.get() && !rc_event_cloned.get().as_ref().read {
            let event_id = event.id.clone();
            spawn_local_scoped(cx, async move {
                debug!("mark_event_as_read: {:?}", event_id.clone());
                mark_event_as_read(cx, event_id).await;
            });
        }
    });

    view! {cx,
        div(ref=event_ref, class="flex flex-col rounded-lg shadow my-4 p-4 w-full bg-gray-100 dark:bg-purple-700 ") {
            div(class="flex flex-row justify-between") {
                div(class="flex flex-row items-center w-full") {
                    div(class="rounded-full h-14 w-14 aspect-square mr-4 bg-gray-300 dark:bg-purple-600 flex items-center justify-center relative") {
                        img(src=event.chain.clone().unwrap().image_url, alt="Event Logo", class="h-12 w-12")
                        div(class=format!("w-2 h-2 bg-red-500 rounded-full absolute top-0 right-0 transition-opacity duration-[3000ms] {}", if rc_event.get().read {"opacity-0"} else {"opacity-100"})) {}
                    }
                    div(class="flex flex-col w-full") {
                        div(class="flex flex-row text-center items-center") {
                            span(class="text-lg font-bold") { (format!("{} {}", event.title.clone(), event.emoji)) }
                            EventBadge(event_type=event_type)
                            div(class="flex-grow") {}
                            span(class="text-sm justify-self-end dark:text-purple-600") { (display_timestamp(event.created_at.clone(), locale.to_string())) }
                        }
                        p(class="text-sm font-bold py-1") { (event.subtitle.clone()) }
                        p(class=format!("text-sm italic {}", if *is_collapsed.get() {"line-clamp-3"} else {""})) { (event.description.clone()) }
                        div(class="flex items-center justify-center") {
                            button(class=format!("flex items-center justify-center text-sm font-bold border rounded-lg mt-2 p-2 \
                                    text-purple-600 border-purple-600 hover:text-primary hover:border-primary {}", if is_clamping {""} else {"hidden"}), on:click=move |_| is_collapsed.set(!*is_collapsed.get())) {
                                (if *is_collapsed.get() {
                                    view! {cx,
                                        div(class="icon-[simple-line-icons--arrow-down]") {}
                                    }
                                } else {
                                    view! {cx,
                                        div(class="icon-[simple-line-icons--arrow-up]") {}
                                    }
                                })
                            }
                        }
                    }
                }
            }
        }
    }
}

#[component]
pub fn WelcomeMessage<G: Html>(cx: Scope) -> View<G> {
    let events_state = use_context::<EventsState>(cx);
    let welcome_wallets = create_memo(cx, || {
        if events_state.welcome_message.get().is_none() {
            return vec![];
        }
        let mut wallets = events_state
            .welcome_message
            .get()
            .as_ref()
            .clone()
            .unwrap()
            .to_vec();
        wallets.sort_by(|a, b| a.name.cmp(&b.name));
        wallets
    });

    create_effect(cx, move || {
        if events_state.welcome_message.get().is_some() && !events_state.events.get().is_empty() {
            *events_state.welcome_message.modify() = None;
        }
    });

    view! {cx,
        (if events_state.welcome_message.get().is_some() {
                view!{cx,
                    div(class="flex flex-col rounded-lg shadow my-4 p-4 w-full bg-gray-100 dark:bg-purple-700") {
                        div(class="flex flex-col rounded-lg shadow my-2 p-2 pl-16 w-full bg-gray-100 dark:bg-purple-700") {
                            h1(class="text-2xl font-bold") { "Welcome on Star Scope" }
                            p(class="text-sm py-1") { "You will receive notifications about:" }
                            ul(class="list-disc list-inside text-sm") {
                                li { "Receiving of funds per transaction or IBC" }
                                li { "End of unstaking perdiod" }
                                li { "End of unbonding period of Osmosis pools" }
                                li { "End of vesting period of Neutron Airdrop" }
                                li { "New governance proposals" }
                                li { "Passes/rejected governance proposals" }
                            }
                            p(class="text-sm font-bold py-1") { "Following wallets are registered" }
                        }
                        div(class="flex flex-row flex-wrap justify-center") {
                            Indexed(
                                iterable=welcome_wallets,
                                view=move |cx,wallet_info| {
                                    view!{cx,
                                        div(class="flex flex-col items-center justify-center rounded-lg shadow my-2 p-2 w-full bg-gray-100 dark:bg-purple-700") {
                                            div(class="flex flex-row items-center w-full") {
                                                div(class="rounded-full h-8 w-8 aspect-square mr-4 bg-gray-300 dark:bg-purple-600 flex items-center justify-center relative") {
                                                    img(src=wallet_info.image_url, alt="Event Logo", class="h-6 w-6")
                                                }
                                                div(class="flex flex-col w-full") {
                                                    div(class="flex flex-row text-center items-center") {
                                                        span(class="text-base font-semibold") { (wallet_info.name.clone()) }
                                                        div(class="flex-grow") {}
                                                    }
                                                    p(class="text-sm py-1") { (wallet_info.address.clone()) }
                                                }
                                            }
                                        }
                                    }
                                }
                            )
                        }
                    }
                }
            } else {
                view!{cx,
                    div()
                }
            })
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
            .filter(|event| {
                let event_type_filter = notifications_state.event_type_filter.get();
                match event_type_filter.as_ref() {
                    None => true,
                    Some(filter) => event.get().event_type() == *filter,
                }
            })
            .filter(|event| {
                let read_status_filter = notifications_state.read_status_filter.get();
                match read_status_filter.as_ref() {
                    ReadStatusFilter::All => true,
                    ReadStatusFilter::Read => event.get().read,
                    ReadStatusFilter::Unread => !event.get().read,
                }
            })
            .filter(|event| {
                let chain_filter = notifications_state.chain_filter.get();
                match chain_filter.as_ref() {
                    None => true,
                    Some(chain) => event.get().chain.clone().unwrap().id == chain.id,
                }
            })
            .filter(|event| {
                let time_filter = notifications_state.time_filter.get();
                match time_filter.as_ref().as_time_range() {
                    None => true,
                    Some((start, end)) => {
                        event.get().created_at.clone().unwrap().seconds > start.timestamp()
                            && event.get().created_at.clone().unwrap().seconds <= end.timestamp()
                    }
                }
            })
            .take(100)
            .cloned()
            .collect::<Vec<_>>()
    });

    view! {cx,
        div(class="flex flex-col") {
            WelcomeMessage {}
            Keyed(
                iterable=events,
                key=|event| event.get().id.clone(),
                view=move |cx,event| {
                    view!{cx,
                        EventComponent(rc_event=event)
                    }
                }
            )
        }
    }
}

#[derive(Debug, Clone)]
pub struct NotificationsState {
    event_type_filter: RcSignal<Option<EventType>>,
    read_status_filter: RcSignal<ReadStatusFilter>,
    chain_filter: RcSignal<Option<grpc::ChainData>>,
    time_filter: RcSignal<TimeFilter>,
    chains: RcSignal<Vec<grpc::ChainData>>,
    locale: RcSignal<String>,
}

impl NotificationsState {
    pub fn new() -> Self {
        let window = web_sys::window().expect("Missing Window");
        let navigator = window.navigator();
        let locale = navigator.language().unwrap_or_else(|| "en-GB".to_string());
        Self {
            event_type_filter: create_rc_signal(None),
            read_status_filter: create_rc_signal(ReadStatusFilter::default()),
            chain_filter: create_rc_signal(None),
            time_filter: create_rc_signal(TimeFilter::default()),
            chains: create_rc_signal(Vec::new()),
            locale: create_rc_signal(locale),
        }
    }

    pub fn reset(&self) {
        self.event_type_filter.set(None);
        self.read_status_filter.set(ReadStatusFilter::default());
        self.chain_filter.set(None);
        self.time_filter.set(TimeFilter::default());
    }

    pub fn add_chains(&self, chains: Vec<grpc::ChainData>) {
        self.chains.set(chains);
    }

    pub fn apply_filter(&self, filter: Option<EventType>) {
        self.event_type_filter.set(filter);
    }

    pub fn has_filter_applied(&self, filter: Option<EventType>) -> bool {
        match filter {
            None => self.event_type_filter.get().is_none(),
            Some(et) => match self.event_type_filter.get().as_ref() {
                None => false,
                Some(f) => *f == et,
            },
        }
    }
}

async fn query_chains(cx: Scope<'_>) {
    let notifications_state = use_context::<NotificationsState>(cx);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(());
    let response = services
        .grpc_client
        .get_event_service()
        .list_chains(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(response) = response {
        notifications_state.add_chains(response.chains);
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

#[derive(Debug, Clone, PartialEq, Sequence)]
pub enum ReadStatusFilter {
    All,
    Read,
    Unread,
}

#[derive(Debug, PartialEq, Eq)]
pub struct ReadStatusFilterError;

impl ReadStatusFilter {
    const QUERY_PARAM: &'static str = "read_status";

    fn get_filter_from_hash() -> Self {
        match get_query_param(ReadStatusFilter::QUERY_PARAM) {
            None => ReadStatusFilter::All,
            Some(param) => match param.as_str() {
                "read" => ReadStatusFilter::Read,
                "unread" => ReadStatusFilter::Unread,
                _ => ReadStatusFilter::All,
            },
        }
    }

    #[allow(dead_code)]
    fn to_hash(&self) -> String {
        if self == &ReadStatusFilter::All {
            "".to_string()
        } else {
            self.to_string()
        }
    }

    #[allow(dead_code)]
    fn set_filter_as_query_param(&self) {
        add_or_update_query_params(ReadStatusFilter::QUERY_PARAM, self.to_hash().as_str());
    }

    fn default() -> Self {
        ReadStatusFilter::get_filter_from_hash()
    }
}

impl Display for ReadStatusFilter {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        match self {
            ReadStatusFilter::All => write!(f, "all"),
            ReadStatusFilter::Read => write!(f, "read"),
            ReadStatusFilter::Unread => write!(f, "unread"),
        }
    }
}

impl FromStr for ReadStatusFilter {
    type Err = ReadStatusFilterError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "all" => Ok(ReadStatusFilter::All),
            "read" => Ok(ReadStatusFilter::Read),
            "unread" => Ok(ReadStatusFilter::Unread),
            _ => Err(ReadStatusFilterError),
        }
    }
}

const DROPDOWN_DIV_CLASS: &str = "relative inline-flex items-center w-full";
const DROPDOWN_ICON_CLASS: &str = "absolute w-full left-0 top-0 h-full flex items-center pl-2 pointer-events-none text-gray-500 dark:text-purple-600";
const DROPDOWN_SELECT_CLASS: &str = "block capitalize w-full md:w-auto pl-8 py-2 rounded border-0 duration-300 hover:bg-sky-400 dark:text-purple-600 dark:bg-purple-700 dark:hover:bg-purple-800";

#[component]
pub fn ReadStatusFilterDropdown<G: Html>(cx: Scope) -> View<G> {
    let notifications_state = use_context::<NotificationsState>(cx);
    let input_ref = create_node_ref(cx);

    let handle_change = |event: Event| {
        let target: HtmlSelectElement = event.target().unwrap().unchecked_into();
        let filter = ReadStatusFilter::from_str(&target.value()).unwrap();
        notifications_state.read_status_filter.set(filter);
        // filter.set_filter_as_query_param();
    };

    let options = View::new_fragment(
        all::<ReadStatusFilter>().map(|f: ReadStatusFilter| {
            let cloned_f: ReadStatusFilter = f.clone();
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

#[component]
pub fn ChainFilterDropdown<G: Html>(cx: Scope) -> View<G> {
    let notifications_state = use_context::<NotificationsState>(cx);
    let input_ref = create_node_ref(cx);
    const ALL_KEY: &str = "all";

    let handle_change = |event: Event| {
        let target: HtmlSelectElement = event.target().unwrap().unchecked_into();
        if target.value() == ALL_KEY {
            notifications_state.chain_filter.set(None);
        } else {
            let chain_id = target.value().parse::<i64>().unwrap();
            if let Some(chain) = notifications_state
                .chains
                .get()
                .iter()
                .find(|c| c.id == chain_id)
            {
                notifications_state.chain_filter.set(Some(chain.clone()));
            } else {
                notifications_state.chain_filter.set(None);
            }
        }
    };

    let chains = create_memo(cx, || {
        let mut chains = notifications_state
            .chains
            .get()
            .iter()
            .map(|c| Some(c.clone()))
            .collect::<Vec<_>>();
        chains.insert(0, None);
        chains
    });

    view! { cx,
        div(class=DROPDOWN_DIV_CLASS) {
            div(class=DROPDOWN_ICON_CLASS) {
                span(class="icon-[system-uicons--cubes]")
            }
            select(ref=input_ref,
                class=DROPDOWN_SELECT_CLASS,
                on:change=handle_change,
            ) {
                Keyed(
                    iterable=chains,
                    key=|chain| {
                        match chain {
                            Some(chain) => chain.id.to_string(),
                            None => ALL_KEY.to_string(),
                        }
                    },
                    view=|cx,chain| {
                        match chain {
                            Some(chain) => {
                                view!{cx,
                                    option(value=chain.id, class="capitalize") { (chain.name) }
                                }
                            },
                            None => {
                                view!{cx,
                                    option(value=ALL_KEY, class="capitalize") { "All" }
                                }
                            }
                        }
                    }
                )
            }
        }
    }
}

#[derive(Debug, Clone, PartialEq, Sequence)]
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
    const QUERY_PARAM: &'static str = "time_filter";

    fn get_filter_from_hash() -> Self {
        match get_query_param(TimeFilter::QUERY_PARAM) {
            None => TimeFilter::All,
            Some(param) => TimeFilter::from_str(param.as_str()).unwrap_or(TimeFilter::All),
        }
    }

    #[allow(dead_code)]
    fn to_hash(&self) -> String {
        if self == &TimeFilter::All {
            "".to_string()
        } else {
            self.to_string()
        }
    }

    #[allow(dead_code)]
    fn set_filter_as_query_param(&self) {
        add_or_update_query_params(TimeFilter::QUERY_PARAM, self.to_hash().as_str());
    }

    fn default() -> Self {
        TimeFilter::get_filter_from_hash()
    }

    fn as_time_range(&self) -> Option<(NaiveDateTime, NaiveDateTime)> {
        let js_date = Date::new_0();
        let milliseconds = js_date.get_time();
        let seconds = (milliseconds / 1000.0) as i64;
        let today = NaiveDateTime::from_timestamp_opt(seconds, 0)
            .unwrap()
            .date()
            .and_hms_opt(0, 0, 0)
            .unwrap();
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

impl Display for TimeFilter {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        match self {
            TimeFilter::All => write!(f, "all"),
            TimeFilter::Today => write!(f, "today"),
            TimeFilter::Yesterday => write!(f, "yesterday"),
            TimeFilter::OneWeek => write!(f, "one week"),
            TimeFilter::OneMonth => write!(f, "one month"),
            TimeFilter::OneYear => write!(f, "one year"),
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
            "one week" => Ok(TimeFilter::OneWeek),
            "one month" => Ok(TimeFilter::OneMonth),
            "one year" => Ok(TimeFilter::OneYear),
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
        let filter = TimeFilter::from_str(&target.value()).unwrap();
        notifications_state.time_filter.set(filter);
        // filter.set_filter_as_query_param();
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
                span(class="icon-[ic--outline-access-time]")
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

async fn query_events(cx: Scope<'_>, event_type: Option<EventType>) {
    debug!("query_events: {:?}", event_type);
    let events_state = use_context::<EventsState>(cx);
    let services = use_context::<Services>(cx);
    let request = services
        .grpc_client
        .create_request(grpc::ListEventsRequest {
            start_time: None,
            end_time: None,
            limit: 0,
            offset: 0,
            event_type: event_type.map(|e| e as i32),
        });
    let response = services
        .grpc_client
        .get_event_service()
        .list_events(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(response) = response {
        events_state.replace_events(response.events, event_type);
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

#[component]
pub async fn Notifications<G: Html>(cx: Scope<'_>) -> View<G> {
    let notifications_state = use_context::<NotificationsState>(cx);

    spawn_local_scoped(cx.to_owned(), async move {
        query_chains(cx.to_owned()).await;
    });

    let event_type_filter =
        create_selector(cx, move || *notifications_state.event_type_filter.get());

    create_effect(cx, move || {
        let event_type = *event_type_filter.get();
        spawn_local_scoped(cx.to_owned(), async move {
            query_events(cx.to_owned(), event_type).await;
        });
    });

    view! {cx,
        div(class="flex flex-col") {
            div(class="hidden lg:flex flex-row justify-between items-center pb-4") {
                h1(class="text-4xl font-bold") { "Notifications" }
                div(class="flex flex-row space-x-4 h-8") {
                    // ReadStatusFilterDropdown {}
                    ChainFilterDropdown {}
                    TimeFilterDropdown {}
                }
            }
            div(class="lg:hidden flex flex-col") {
                h1(class="text-4xl font-bold pb-4") { "Notifications" }
                div(class="flex flex-wrap") {
                    // div(class="w-full sm:w-auto flex-shrink-0 flex-grow-0 mb-4 sm:mb-0 sm:mr-4") {
                    //     ReadStatusFilterDropdown {}
                    // }
                    div(class="w-full sm:w-auto flex-shrink-0 flex-grow-0 mb-4 sm:mb-0 sm:mr-4") {
                        ChainFilterDropdown {}
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
