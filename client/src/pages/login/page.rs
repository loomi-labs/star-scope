use log::debug;
use serde::{Deserialize, Serialize};
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;
use wasm_bindgen::prelude::*;

use crate::{AppState, AuthState, InfoLevel, Services};
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
    } else if js_result.result.is_empty() {
        return Err("Keplr login failed".to_string());
    }
    Ok(js_result.result)
}

#[component]
pub async fn Login<G: Html>(cx: Scope<'_>) -> View<G> {
    let app_state = use_context::<AppState>(cx);

    view!(cx,
        div(class="h-screen w-screen bg-gray-100 dark:bg-purple-900 flex flex-col justify-center py-12 sm:px-6 lg:px-8") {
            div(class="sm:mx-auto sm:w-full sm:max-w-md") {
                img(class="mx-auto h-12 w-auto", src="https://tailwindui.com/img/logos/workflow-mark-indigo-600.svg", alt="Workflow")
                h2(class="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white") {
                    "Login with Keplr"
                }
            }
            div(class="mt-8 sm:mx-auto sm:w-full sm:max-w-md") {
                div(class="bg-white dark:bg-purple-700 py-8 px-4 shadow sm:rounded-lg sm:px-10") {
                    div(class="space-y-6") {
                        div { /* Submit button */
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
                                                        create_message(cx, "Login success", format!("Logged in successfully"), InfoLevel::Info);
                                                    }
                                                    Err(status) => create_message(cx, "Login failed", format!("Login failed with status: {}", status), InfoLevel::Error),
                                                }
                                            }
                                            Err(status) => create_message(cx, "Login failed", format!("Login failed with status: {}", status), InfoLevel::Error),
                                        }
                                    };
                                });
                            }, class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500") {
                                "Keplr Login"
                            }
                        }
                    }
                }
            }
        }
    )
}
