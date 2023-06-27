use crate::components::messages::create_message;
use crate::types::protobuf::grpc::{DiscordLoginRequest, TelegramLoginRequest};
use crate::{AppRoutes, AppState, AuthState, InfoLevel};
use log::debug;
use sycamore::prelude::{use_context, Scope};
use sycamore_router::navigate;
use urlencoding::decode;
use wasm_bindgen::JsValue;

#[allow(dead_code)]
pub fn add_or_update_query_params(key: &str, value: &str) {
    let window = web_sys::window().unwrap();
    let history = window.history().unwrap();
    let document = window.document().unwrap();
    let url = document.url().unwrap();
    let segments = url.split('?').collect::<Vec<_>>();
    let stripped_url = segments[0];
    let mut query_params = querystring::querify("");
    if segments.len() > 1 {
        query_params = querystring::querify(segments[1]);
    }
    let index = query_params.iter().position(|(k, _)| *k == key);
    if let Some(index) = index {
        query_params.remove(index);
    }
    if !value.is_empty() {
        query_params.push((key, value));
    }
    let new_query_params = querystring::stringify(vec![(key, value)])
        .trim_end_matches('&')
        .to_string();
    let new_url_string = format!("{}?{}", stripped_url, new_query_params);
    history
        .replace_state_with_url(&JsValue::null(), "", Some(&new_url_string))
        .unwrap();
}

pub fn clean_query_params() {
    let window = web_sys::window().unwrap();
    let history = window.history().unwrap();
    let document = window.document().unwrap();
    let url = document.url().unwrap();
    let segments = url.split('?').collect::<Vec<_>>();
    let stripped_url = segments[0];
    history
        .replace_state_with_url(&JsValue::null(), "", Some(stripped_url))
        .unwrap();
}

pub fn get_query_param(key: &str) -> Option<String> {
    let window = web_sys::window().unwrap();
    let document = window.document().unwrap();
    let url = document.url().unwrap();
    let segments = url.split('?').collect::<Vec<_>>();
    if segments.len() < 2 {
        return None;
    }
    let query_params = querystring::querify(segments[1]);
    let index = query_params.iter().position(|(k, _)| *k == key);
    index.map(|index| query_params[index].1.to_string())
}

pub fn get_query_params() -> Vec<(String, String)> {
    let location = web_sys::window().unwrap().location();
    let search: Result<String, JsValue> = location.search();
    let mut params = Vec::new();
    for s in search.unwrap().trim_start_matches('?').split('&') {
        if s.is_empty() {
            continue;
        }
        let mut kv = s.split('=');
        let k = kv.next().unwrap();
        let v = kv.next().unwrap();
        params.push((k.to_string(), v.to_string()));
    }
    params
}

pub fn has_telegram_login_query_params() -> bool {
    get_query_param("hash").is_some()
}

pub fn has_discord_login_query_params() -> bool {
    get_query_param("code").is_some()
}

pub fn get_discord_login_data() -> Option<DiscordLoginRequest> {
    get_query_param("code").map(|code| DiscordLoginRequest { code })
}

pub fn get_telegram_login_data() -> Option<TelegramLoginRequest> {
    let query_params = get_query_params();
    let hash = query_params
        .iter()
        .find(|params: &&(String, String)| params.0 == "hash")
        .map(|params: &(String, String)| params.1.clone())?;
    let mut data = query_params
        .iter()
        .filter(|params: &&(String, String)| params.0 != "hash")
        .map(|params: &(String, String)| format!("{}={}", params.0, params.1))
        .collect::<Vec<_>>();
    data.sort();
    let data_str = data.join("\n");
    let data_str_decoded = decode(data_str.as_str()).expect("UTF-8");

    Some(TelegramLoginRequest {
        data_str: data_str_decoded.to_string(),
        hash,
    })
}

pub fn safe_navigate(cx: Scope, route: AppRoutes) {
    let app_state = use_context::<AppState>(cx);
    if app_state.route.get_untracked().is_some_and(|r| r != route) {
        navigate(route.to_string().as_str());
    }
}

pub fn navigate_launch_app(cx: Scope) {
    let app_state = use_context::<AppState>(cx);
    if *app_state.auth_state.get_untracked() == AuthState::LoggedIn {
        if let Some(user) = app_state.user.get_untracked().as_ref() {
            if user.is_setup_complete {
                debug!("Redirect to notifications");
                safe_navigate(cx, AppRoutes::Notifications)
            } else {
                debug!("Redirect to setup");
                safe_navigate(cx, AppRoutes::Setup)
            }
        } else {
            create_message(cx, "User not found", "User status unknown", InfoLevel::Error);
        }
    } else {
        safe_navigate(cx, AppRoutes::Login)
    }
}
