use std::result;

use log::debug;
use serde::{Deserialize, Serialize};
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;
use wasm_bindgen::prelude::*;

use crate::{AppRoutes, AppState, AuthState, InfoLevel, Services};
use crate::components::messages::create_message;

#[derive(Serialize, Deserialize)]
pub struct JsResult {
    pub result: String,
    pub error: String,
}

#[wasm_bindgen(module = "/src/pages/login/keplr_login.js")]
extern "C" {
    async fn keplr_login() -> JsValue;
}

async fn keplr_login_wrapper() -> Result<String, String> {
    let login_result = keplr_login().await;
    let js_result: JsResult = serde_wasm_bindgen::from_value(login_result).unwrap();
    if !js_result.error.is_empty() {
        return Err(js_result.error);
    } else if js_result.result == "" {
        return Err("Keplr login failed".to_string());
    }
    Ok(js_result.result)
}

#[component]
pub async fn Login<G: Html>(cx: Scope<'_>) -> View<G> {
    let app_state = use_context::<AppState>(cx);

    view!(cx,
        button(on:click=move |_| {
            spawn_local_scoped(cx, async move {
                if *app_state.auth_state.get() == AuthState::LoggedOut {
                    debug!("Attempt to login");
                    match keplr_login_wrapper().await {
                        Ok(result) => {
                            let response = use_context::<Services>(cx).auth_manager.clone().login(result.clone()).await;
                            match response {
                                Ok(_) => {
                                    let mut auth_state = use_context::<AppState>(cx).auth_state.modify();
                                    *auth_state = AuthState::LoggedIn;
                                    create_message(cx, "Login success", &format!("Logged in successfully"), InfoLevel::Info);
                                }
                                Err(status) => create_message(cx, "Login failed", &format!("Login failed with status: {}", status), InfoLevel::Error),
                            }
                        }
                        Err(status) => create_message(cx, "Login failed", &format!("Login failed with status: {}", status), InfoLevel::Error),
                    }
                };
            });
        }) { "Login" }
        a(href=AppRoutes::Home) { "Home" }
        a(href="keplr", rel="external") { "Keplr" }
    )
}
