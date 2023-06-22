#![allow(non_snake_case)]

use std::cmp::Ordering;
use std::collections::HashMap;
use std::fmt::Display;

use log::{debug, error};
use log::Level;
use prost_types::Timestamp;
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;
use sycamore::suspense::Suspense;
use sycamore_router::{HistoryIntegration, Route, Router};
use uuid::Uuid;

use crate::components::layout::LayoutWrapper;
use crate::components::messages::{create_error_msg_from_status, create_message, MessageOverlay};
use crate::pages::communication::page::Communication;
use crate::pages::home::page::Home;
use crate::pages::login::page::Login;
use crate::pages::notifications::page::{Notifications, NotificationsState};
use crate::pages::settings::page::Settings;
use crate::services::auth::AuthService;
use crate::services::grpc::GrpcClient;
use crate::types::protobuf::event::EventType;
use crate::types::protobuf::grpc::{Event, EventsCount, User, WalletInfo};
use crate::utils::url::safe_navigate;

mod components;
mod config;
mod pages;
mod services;
mod types;
mod utils;

#[derive(Route, Debug, Clone, Copy, PartialEq)]
pub enum AppRoutes {
    #[to("/")]
    Home,
    #[to("/notifications")]
    Notifications,
    #[to("/communication")]
    Communication,
    #[to("/settings")]
    Settings,
    #[to("/login")]
    Login,
    #[not_found]
    NotFound,
}

impl AppRoutes {
    fn needs_login(&self) -> bool {
        match self {
            AppRoutes::Home => false,
            AppRoutes::Notifications => true,
            AppRoutes::Communication => true,
            AppRoutes::Settings => true,
            AppRoutes::Login => false,
            AppRoutes::NotFound => false,
        }
    }
}

impl ToString for AppRoutes {
    fn to_string(&self) -> String {
        match self {
            AppRoutes::Home => "/".to_string(),
            AppRoutes::Notifications => "/notifications".to_string(),
            AppRoutes::Communication => "/communication".to_string(),
            AppRoutes::Settings => "/settings".to_string(),
            AppRoutes::Login => "/login".to_string(),
            AppRoutes::NotFound => "/404".to_string(),
        }
    }
}

#[derive(Debug, Clone)]
pub struct Services {
    pub grpc_client: GrpcClient,
    pub auth_manager: AuthService,
}

impl Services {
    pub fn new() -> Self {
        Self {
            grpc_client: GrpcClient::default(),
            auth_manager: AuthService::default(),
        }
    }
}

impl Default for Services {
    fn default() -> Self {
        Self::new()
    }
}

#[derive(Debug, Clone, PartialEq)]
pub enum AuthState {
    LoggedOut,
    LoggedIn,
    LoggingIn,
}

impl Display for AuthState {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            AuthState::LoggedOut => write!(f, "LoggedOut"),
            AuthState::LoggedIn => write!(f, "LoggedIn"),
            AuthState::LoggingIn => write!(f, "LoggingIn"),
        }
    }
}

#[repr(usize)]
#[derive(Debug, Clone, Eq, PartialEq)]
pub enum InfoLevel {
    Info = 1,
    Success,
    Error,
}

#[derive(Debug, Clone, PartialEq)]
pub struct InfoMsg {
    pub id: Uuid,
    pub title: String,
    pub message: String,
    pub level: InfoLevel,
    pub created_at: f64,
}

#[derive(Debug, Clone)]
pub struct AppState {
    auth_service: AuthService,
    pub auth_state: RcSignal<AuthState>,
    pub route: RcSignal<Option<AppRoutes>>,
    pub messages: RcSignal<Vec<RcSignal<InfoMsg>>>,
    pub user: RcSignal<Option<User>>,
    pub is_dialog_open: RcSignal<bool>,
}

impl Default for AppState {
    fn default() -> Self {
        Self::new(AuthService::default())
    }
}

impl AppState {
    pub fn new(auth_service: AuthService) -> Self {
        let auth_state = match auth_service.is_jwt_valid() {
            true => AuthState::LoggedIn,
            false => AuthState::LoggedOut,
        };
        Self {
            auth_service,
            auth_state: create_rc_signal(auth_state),
            route: create_rc_signal(None),
            messages: create_rc_signal(vec![]),
            user: create_rc_signal(None),
            is_dialog_open: create_rc_signal(false),
        }
    }

    pub fn add_message(&self, title: String, message: String, level: InfoLevel) -> Uuid {
        let uuid = Uuid::new_v4();
        let created_at = js_sys::Date::now();
        self.messages.modify().push(create_rc_signal(InfoMsg {
            id: uuid,
            title,
            message,
            level,
            created_at,
        }));
        uuid
    }

    pub fn remove_message(&self, id: Uuid) {
        self.messages.modify().retain(|m| m.get().id != id);
    }

    pub fn logout(&self) {
        self.auth_service.logout();
        self.user.set(None);
        self.auth_state.set(AuthState::LoggedOut);
    }

    pub fn set_showing_dialog(&self, is_open: bool) {
        self.is_dialog_open.set(is_open);
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct EventsState {
    pub events: RcSignal<Vec<RcSignal<Event>>>,
    pub event_count_map: RcSignal<HashMap<EventType, EventsCount>>,
    pub welcome_message: RcSignal<Option<Vec<WalletInfo>>>,
}

impl Default for EventsState {
    fn default() -> Self {
        Self::new()
    }
}

fn compare_timestamps(a: &Option<Timestamp>, b: &Option<Timestamp>) -> Ordering {
    match (a, b) {
        (Some(a), Some(b)) => a.seconds.cmp(&b.seconds).then(a.nanos.cmp(&b.nanos)),
        (Some(_), None) => Ordering::Greater,
        (None, Some(_)) => Ordering::Less,
        (None, None) => Ordering::Equal,
    }
}

impl EventsState {
    pub fn new() -> Self {
        Self {
            events: create_rc_signal(vec![]),
            event_count_map: create_rc_signal(HashMap::new()),
            welcome_message: create_rc_signal(None),
        }
    }

    pub fn reset(&self) {
        self.events.set(vec![]);
    }

    pub fn replace_events(&self, new_events: Vec<Event>, _event_type: Option<EventType>) {
        let mut events = self.events.modify();
        for e in new_events.iter() {
            events.retain(|e2| e2.get().id != e.id);
            events.insert(0, create_rc_signal(e.clone()));
        }
        events.sort_by(|a, b| compare_timestamps(&a.get().notify_at, &b.get().notify_at));
        events.reverse();
        *self.events.modify() = events.clone();
    }

    pub fn update_event_count(&self, events_count: Vec<EventsCount>) {
        let mut event_count_map = self.event_count_map.modify();
        for e in events_count {
            event_count_map.insert(e.event_type(), e);
        }
    }

    pub fn mark_as_read(&self, event_id: String) {
        let events = self.events.modify();
        if let Some(index) = events.iter().position(|e| e.get().id == event_id) {
            let mut event = events[index].get().as_ref().clone();
            event.read = true;
            *events[index].modify() = event.clone();
            let event_type = event.event_type();
            let mut event_count_map = self.event_count_map.modify();
            if let Some(count) = event_count_map.get_mut(&event_type) {
                count.unread_count -= 1;
            }
            *self.event_count_map.modify() = event_count_map.clone();
        }
        *self.events.modify() = events.clone();
    }

    pub fn has_events(&self) -> bool {
        self.event_count_map.get().iter().any(|(_, v)| v.count > 0)
    }

    pub fn add_welcome_message(&self, message: Vec<WalletInfo>) {
        self.welcome_message.set(Some(message));
    }
}

fn start_jwt_refresh_timer(cx: Scope) {
    spawn_local_scoped(cx, async move {
        gloo_timers::future::TimeoutFuture::new(1000 * 60).await;
        let auth_client = AuthService::new();
        debug!("is_jwt_valid: {}", auth_client.is_jwt_valid());
        if auth_client.is_jwt_about_to_expire() {
            auth_client.refresh_access_token().await;
        }
        if auth_client.is_jwt_valid() {
            start_jwt_refresh_timer(cx.to_owned());
        } else {
            debug!("JWT is not valid anymore");
            let app_state = use_context::<AppState>(cx);
            app_state.logout();
        }
    });
}

fn has_access_permission(auth_service: &AuthService, route: &AppRoutes) -> bool {
    let is_admin = auth_service.is_admin();
    let is_user = auth_service.is_user();
    match route {
        AppRoutes::Home => true,
        AppRoutes::Notifications => is_user || is_admin,
        AppRoutes::Communication => is_user || is_admin,
        AppRoutes::Settings => is_user || is_admin,
        AppRoutes::Login => true,
        AppRoutes::NotFound => true,
    }
}

fn activate_view<G: Html>(cx: Scope, route: &AppRoutes) -> View<G> {
    debug!("Route changed to: {:?}", route);
    let app_state = use_context::<AppState>(cx);
    let services = use_context::<Services>(cx);
    if has_access_permission(&services.auth_manager, route) {
        app_state.route.set(Some(*route));
        match route {
            AppRoutes::Home => view!(cx, Home {}),
            AppRoutes::Notifications => view!(cx, Notifications {}),
            AppRoutes::Communication => view!(cx, Communication {}),
            AppRoutes::Settings => view!(cx, Settings {}),
            AppRoutes::Login => view!(cx, Login {}),
            AppRoutes::NotFound => view! { cx, "404 Not Found"},
        }
    } else {
        app_state.route.set(Some(AppRoutes::Login));
        create_message(
            cx,
            "Access denied".to_string(),
            "Please login to access this page",
            InfoLevel::Info,
        );
        Login(cx)
    }
}

async fn query_user_info(cx: Scope<'_>) {
    let app_state = use_context::<AppState>(cx);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(());
    let response = services
        .grpc_client
        .get_user_service()
        .get_user(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(user) = response {
        *app_state.user.modify() = Some(user);
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

fn subscribe_to_events(cx: Scope) {
    spawn_local_scoped(cx, async move {
        let app_state = use_context::<AppState>(cx);
        match app_state.auth_state.get().as_ref() {
            AuthState::LoggedOut => return,
            AuthState::LoggedIn => {}
            AuthState::LoggingIn => return,
        }

        // let overview_state = use_context::<EventsState>(cx);
        let services = use_context::<Services>(cx);
        let result = services
            .grpc_client
            .get_event_service()
            .event_stream(services.grpc_client.create_request(()))
            .await
            .map(|res| res.into_inner());
        match result {
            Ok(mut event_stream) => {
                loop {
                    match event_stream.message().await {
                        Ok(Some(response)) => {
                            if response.event_type.is_some() {
                                debug!("Received {:?} event", response.clone().event_type());
                                query_events_count(cx).await;
                                // overview_state.add_event_count(response);
                            } else {
                                debug!("Received keep alive event");
                            }
                            match app_state.auth_state.get().as_ref() {
                                AuthState::LoggedOut => break,
                                AuthState::LoggedIn => {}
                                AuthState::LoggingIn => break,
                            }
                        }
                        Ok(None) => {
                            // No more events, exit the loop
                            break;
                        }
                        Err(err) => {
                            create_error_msg_from_status(cx, err);
                            gloo_timers::future::TimeoutFuture::new(1000 * 5).await;
                            subscribe_to_events(cx.to_owned());
                        }
                    }
                }
            }
            Err(err) => {
                error!("Error while subscribing to events: {:?}", err);
                gloo_timers::future::TimeoutFuture::new(1000 * 5).await;
                subscribe_to_events(cx.to_owned());
            }
        }
    });
}

async fn query_welcome_message(cx: Scope<'_>) {
    let events_state = use_context::<EventsState>(cx);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(());
    let response = services
        .grpc_client
        .get_event_service()
        .get_welcome_message(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(response) = response {
        events_state.add_welcome_message(response.wallets)
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

async fn query_events_count(cx: Scope<'_>) {
    let events_state = use_context::<EventsState>(cx);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(());
    let response = services
        .grpc_client
        .get_event_service()
        .list_events_count(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(response) = response {
        events_state.update_event_count(response.counters);
        if !events_state.has_events() {
            query_welcome_message(cx.to_owned()).await;
        }
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

async fn login_by_query_params(cx: Scope<'_>) {
    let services = use_context::<Services>(cx);
    if services.auth_manager.clone().has_login_query_params() {
        let response = use_context::<Services>(cx).auth_manager.clone().login_with_query_params().await;
        match response {
            Ok(_) => {
                let mut auth_state = use_context::<AppState>(cx).auth_state.modify();
                *auth_state = AuthState::LoggedIn;
            }
            Err(status) => {
                create_error_msg_from_status(cx, status);
                safe_navigate(cx, AppRoutes::Home);
            }
        }
    }
}

fn execute_logged_out_fns(cx: Scope<'_>) {
    debug!("Logged out");
    let services = use_context::<Services>(cx);
    spawn_local_scoped(cx, async move {
        if services.auth_manager.clone().has_login_query_params() {
            login_by_query_params(cx.to_owned()).await;
        } else {
            let app_state = use_context::<AppState>(cx);
            if app_state.route.get_untracked().as_ref().is_some_and(|route| route.needs_login()) {
                safe_navigate(cx, AppRoutes::Home);
            }
        }
    });
}

fn execute_logged_in_fns(cx: Scope<'_>) {
    debug!("Logged in");
    let event_state = use_context::<EventsState>(cx);
    let notifications_state = use_context::<NotificationsState>(cx);
    event_state.reset();
    notifications_state.reset();
    spawn_local_scoped(cx, async move {
        let app_state = use_context::<AppState>(cx);
        if app_state.route.get_untracked().as_ref().is_some_and(|route| !route.needs_login()) {
            debug!("Redirect to notifications");
            safe_navigate(cx, AppRoutes::Notifications)
        }
        query_user_info(cx).await;
        query_events_count(cx).await;
        subscribe_to_events(cx);
    });
}

#[component]
pub async fn App<G: Html>(cx: Scope<'_>) -> View<G> {
    let services = Services::new();
    let app_state = AppState::new(services.auth_manager.clone());

    provide_context(cx, services);
    provide_context(cx, app_state.clone());
    provide_context(cx, EventsState::new());
    provide_context(cx, NotificationsState::new());

    start_jwt_refresh_timer(cx.to_owned());

    view! {cx,
        div(class="relative") {
            div(class="bg-white dark:bg-d-bg text-black dark:text-white antialiased") {
                MessageOverlay {}
                (if *app_state.is_dialog_open.get() {
                    let app_state = use_context::<AppState>(cx);
                    view!{cx,
                        div(class="fixed inset-0 bg-black opacity-50 z-50", on:click=move |_| app_state.is_dialog_open.set(false)) {}
                    }
                } else {
                    view!{cx,
                        div {}
                    }
                })
                Router(
                    integration=HistoryIntegration::new(),
                    view=|cx, route: &ReadSignal<AppRoutes>| {
                        let auth_state_changed = create_selector(cx, move || app_state.auth_state.get().as_ref().clone());
                        create_effect(cx, move || {
                            match auth_state_changed.get().as_ref() {
                                AuthState::LoggedOut => execute_logged_out_fns(cx),
                                AuthState::LoggedIn => execute_logged_in_fns(cx),
                                AuthState::LoggingIn => {}
                            }
                        });
                        let has_layout_wrapper = create_selector(cx, move || route.get().as_ref().needs_login());
                        view! {cx,
                            (if *has_layout_wrapper.get() {
                                view! {cx,
                                    LayoutWrapper{ (activate_view(cx, route.get().as_ref())) }
                                }
                            } else {
                                activate_view(cx, route.get().as_ref())
                            })
                        }
                    }
                )
            }
        }
    }
}

fn main() {
    console_error_panic_hook::set_once();
    console_log::init_with_level(Level::Debug).unwrap();
    debug!("Console log level set to debug");

    sycamore::render(|cx| {
        view! { cx,
            Suspense(fallback=components::loading::LoadingSpinner(cx)) {
                App {}
            }
        }
    });
}
