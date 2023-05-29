#![allow(non_snake_case)]

use std::collections::HashMap;
use std::fmt::Display;

use log::{debug, error};
use log::Level;
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;
use sycamore::suspense::Suspense;
use sycamore_router::{HistoryIntegration, navigate, Route, Router};
use tonic::{Status, Streaming};
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
use crate::types::types::grpc::{Event, EventsCount, EventType, User};
use crate::utils::url::safe_navigate;

mod components;
mod config;
mod pages;
mod services;
mod utils;
mod types;

#[derive(Route, Debug, Clone, PartialEq)]
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
    Warning,
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
    pub route: RcSignal<AppRoutes>,
    pub messages: RcSignal<Vec<RcSignal<InfoMsg>>>,
    pub user: RcSignal<Option<User>>,
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
            route: create_rc_signal(AppRoutes::default()),
            messages: create_rc_signal(vec![]),
            user: create_rc_signal(None),
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
}


#[derive(Debug, Clone, PartialEq)]
pub struct EventsState {
    pub events: RcSignal<Vec<Event>>,
    pub event_count_map: RcSignal<HashMap<EventType, EventsCount>>,
}

impl Default for EventsState {
    fn default() -> Self {
        Self::new()
    }
}


impl EventsState {
    pub fn new() -> Self {
        Self {
            events: create_rc_signal(vec![]),
            event_count_map: create_rc_signal(HashMap::new()),
        }
    }

    pub fn reset(&self) {
        self.events.set(vec![]);
    }

    pub fn add_events(&self, new_events: Vec<Event>) {
        if new_events.len() == 0 {
            return;
        }
        let mut events = self.events.modify();
        let mut event_count_map = self.event_count_map.modify();
        for e in new_events {
            events.insert(0, e.clone());
            let count = event_count_map.get(&e.clone().event_type());
            if let Some(count) = count {
                let mut new_count = count.clone();
                new_count.count += 1;
                new_count.unread_count += 1;
                event_count_map.insert(e.event_type().clone(), new_count);
                continue;
            } else {
                event_count_map.insert(
                    e.event_type().clone(),
                    EventsCount {
                        event_type: e.event_type().clone() as i32,
                        count: 1,
                        unread_count: 1,
                    },
                );
            }
        }
        *self.events.modify() = events.clone();
        *self.event_count_map.modify() = event_count_map.clone();
    }

    pub fn replace_events(&self, new_events: Vec<Event>, event_type: Option<EventType>) {
        let mut events = self.events.modify();
        let f = new_events
            .iter()
            .filter(|e| {
                if e.id == 0 {
                    if let Some(event_type) = event_type {
                        return e.event_type != event_type as i32;
                    }
                }
                true
            });
        for e in f {
            if !events.contains(&e) {
                events.insert(0, e.clone());
            }
        }
        *self.events.modify() = events.clone();
    }

    pub fn update_event_count(&self, events_count: Vec<EventsCount>) {
        let mut event_count_map = self.event_count_map.modify();
        for e in events_count {
            event_count_map.insert(e.event_type(), e);
        }
        *self.event_count_map.modify() = event_count_map.clone();
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
        app_state.route.set(route.clone());
        match route {
            AppRoutes::Home => view!(cx, Home {}),
            AppRoutes::Notifications => view!(cx, LayoutWrapper{Notifications {}}),
            AppRoutes::Communication => view!(cx, LayoutWrapper{Communication {}}),
            AppRoutes::Settings => view!(cx, LayoutWrapper{Settings {}}),
            AppRoutes::Login => Login(cx),
            AppRoutes::NotFound => view! { cx, "404 Not Found"},
        }
    } else {
        app_state.route.set(AppRoutes::Login);
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
        let overview_state = use_context::<EventsState>(cx);
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
                            debug!("Received {:?} events", response.events.len());
                            overview_state.add_events(response.events);
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

async fn query_events_count(cx: Scope<'_>) {
    let events_state = use_context::<EventsState>(cx);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request({});
    let response = services
        .grpc_client
        .get_event_service()
        .list_events_count(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(response) = response {
        events_state.update_event_count(response.counters);
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

#[component]
pub async fn App<G: Html>(cx: Scope<'_>) -> View<G> {
    let services = Services::new();
    let app_state = AppState::new(services.auth_manager.clone());

    provide_context(cx, services.clone());
    provide_context(cx, app_state);
    provide_context(cx, EventsState::new());
    provide_context(cx, NotificationsState::new());

    start_jwt_refresh_timer(cx.to_owned());

    view! {cx,
        div(class="bg-white dark:bg-d-bg text-black dark:text-white antialiased") {
            MessageOverlay {}
            Router(
                integration=HistoryIntegration::new(),
                view=|cx, route: &ReadSignal<AppRoutes>| {
                    debug!("Router: create_effect");
                    create_effect(cx, move || {
                        let app_state = use_context::<AppState>(cx);
                        let auth_state = app_state.auth_state.get();
                        debug!("Auth state changed: {}", auth_state);
                        match auth_state.as_ref() {
                            AuthState::LoggedOut => safe_navigate(cx, AppRoutes::Home),
                            AuthState::LoggedIn => {
                                let event_state = use_context::<EventsState>(cx);
                                let notifications_state = use_context::<NotificationsState>(cx);
                                event_state.reset();
                                notifications_state.reset();
                                spawn_local_scoped(cx, async move {
                                    query_user_info(cx).await;
                                    query_events_count(cx).await;
                                    subscribe_to_events(cx);
                                });
                                navigate(AppRoutes::Notifications.to_string().as_str())
                            },
                            AuthState::LoggingIn => {}
                        }
                    });
                    view! {cx, (
                            activate_view(cx, route.get().as_ref())
                        )
                    }
                }
            )
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
