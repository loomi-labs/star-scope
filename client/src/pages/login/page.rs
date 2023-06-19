use serde::{Deserialize, Serialize};
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;
use urlencoding::encode;
use wasm_bindgen::prelude::*;

use crate::components::messages::{create_error_msg_from_status, create_message};
use crate::config::keys;
use crate::{AppState, AuthState, InfoLevel, Services};


#[component]
fn KeplrSvg<G: Html>(cx: Scope) -> View<G> {
    view! (cx,
        svg(
            width="24",
            height="24",
            viewBox="0 0 42 42",
            fill="none",
            xmlns="http://www.w3.org/2000/svg",
        ) {
            g(clip-path="url(#clip0_425_5107)") {
                path(
                    d="M32.4545 0H9.54545C4.27365 0 0 4.27365 0 9.54545V32.4545C0 37.7264 4.27365 42 9.54545 42H32.4545C37.7264 42 42 37.7264 42 32.4545V9.54545C42 4.27365 37.7264 0 32.4545 0Z",
                    fill="url(#paint0_linear_425_5107)",
                )
                path(
                    d="M32.4545 0H9.54545C4.27365 0 0 4.27365 0 9.54545V32.4545C0 37.7264 4.27365 42 9.54545 42H32.4545C37.7264 42 42 37.7264 42 32.4545V9.54545C42 4.27365 37.7264 0 32.4545 0Z",
                    fill="url(#paint1_radial_425_5107)",
                )
                path(
                    d="M32.4545 0H9.54545C4.27365 0 0 4.27365 0 9.54545V32.4545C0 37.7264 4.27365 42 9.54545 42H32.4545C37.7264 42 42 37.7264 42 32.4545V9.54545C42 4.27365 37.7264 0 32.4545 0Z",
                    fill="url(#paint2_radial_425_5107)",
                )
                path(
                    d="M32.4545 0H9.54545C4.27365 0 0 4.27365 0 9.54545V32.4545C0 37.7264 4.27365 42 9.54545 42H32.4545C37.7264 42 42 37.7264 42 32.4545V9.54545C42 4.27365 37.7264 0 32.4545 0Z",
                    fill="url(#paint3_radial_425_5107)",
                )
                path(
                    d="M17.2526 32.2614V22.5192L26.7185 32.2614H31.9849V32.0079L21.0964 20.9122L31.1469 10.3857V10.2614H25.8464L17.2526 19.5635V10.2614H12.9849V32.2614H17.2526Z",
                    fill="white",
                )
            }
            defs {
                linearGradient(
                    id="paint0_linear_425_5107",
                    x1="21",
                    y1="0",
                    x2="21",
                    y2="42",
                    gradientUnits="userSpaceOnUse",
                ) {
                    stop(stop-color="#1FD1FF")
                    stop(offset="1", stop-color="#1BB8FF")
                }
                radialGradient(
                    id="paint1_radial_425_5107",
                    cx="0",
                    cy="0",
                    r="1",
                    gradientUnits="userSpaceOnUse",
                    gradientTransform="translate(2.00623 40.4086) rotate(-45.1556) scale(67.3547 68.3624)",
                ) {
                    stop(stop-color="#232DE3")
                    stop(offset="1", stop-color="#232DE3", stop-opacity="0")
                }
                radialGradient(
                    id="paint2_radial_425_5107",
                    cx="0",
                    cy="0",
                    r="1",
                    gradientUnits="userSpaceOnUse",
                    gradientTransform="translate(39.7379 41.7602) rotate(-138.45) scale(42.1137 64.2116)",
                ) {
                    stop(stop-color="#8B4DFF")
                    stop(offset="1", stop-color="#8B4DFF", stop-opacity="0")
                }
                radialGradient(
                    id="paint3_radial_425_5107",
                    cx="0",
                    cy="0",
                    r="1",
                    gradientUnits="userSpaceOnUse",
                    gradientTransform="translate(20.6501 0.311498) rotate(90) scale(33.1135 80.3423)",
                ) {
                    stop(stop-color="#24D5FF")
                    stop(offset="1", stop-color="#1BB8FF", stop-opacity="0")
                }
                clipPath(id="clip0_425_5107") {
                    rect(width="42", height="42", fill="white")
                }
            }
        }
    )
}

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

#[wasm_bindgen(module = "/src/pages/login/wallet_connect_login.js")]
extern "C" {
    fn wallet_connect_login(url: String) -> JsValue;

    fn isMobile(url: String) -> JsValue;
}

#[allow(dead_code)]
fn wallet_connect_login_wrapper() -> Result<String, String> {
    let login_result = wallet_connect_login(keys::WEB_APP_URL.to_string());
    let js_result = serde_wasm_bindgen::from_value(login_result).unwrap_or_else(|_| JsResult {
        result: "".to_string(),
        error: "Wallet connect login failed".to_string(),
    });
    if !js_result.error.is_empty() {
        return Err(js_result.error);
    } else if js_result.result.is_empty() {
        return Err("Wallet Connect login failed".to_string());
    }
    Ok(js_result.result)
}

fn is_mobile() -> bool {
    let result = isMobile(keys::WEB_APP_URL.to_string());
    serde_wasm_bindgen::from_value(result).unwrap_or(false)
}

#[component]
pub async fn Login<G: Html>(cx: Scope<'_>) -> View<G> {
    let app_state = use_context::<AppState>(cx);
    let class_button = "flex items-center justify-center py-2 px-4 rounded-md shadow-sm text-sm \
    font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500";

    let is_mobile = is_mobile();

    let discord_login_url = format!(
        "https://discord.com/api/oauth2/authorize?client_id={}&redirect_uri={}&response_type=code&scope=identify",
        keys::DISCORD_CLIENT_ID,
        encode(keys::WEB_APP_URL)
    );

    let color = keys::WHITE_COLOR;
    view!(cx,
        div(class="h-screen w-screen bg-gray-100 dark:bg-d-bg flex flex-col justify-center py-12 sm:px-6 lg:px-8") {
            div(class="sm:mx-auto sm:w-full sm:max-w-md flex flex-col justify-center items-center") {
                svg(xmlns="http://www.w3.org/2000/svg", width="200", height="200", viewBox="0 0 563 547.33") {
                    path(fill=color, d="M282.22,24c3.39,42.11,5.12,84.23,6.35,126.34s1.69,84.22,1.64,126.33-.6,84.22-1.77,126.33-2.95,84.23-6.22,126.34c-3.26-42.11-5-84.23-6.22-126.34s-1.73-84.22-1.76-126.33.45-84.22,1.64-126.33S278.84,66.07,282.22,24Z")
                    path(fill=color, d="M538.55,275.5c-42.11,3.35-84.22,5.05-126.33,6.24s-84.23,1.63-126.34,1.54-84.22-.66-126.33-1.87-84.22-3-126.33-6.32c42.12-3.23,84.23-4.9,126.34-6.11s84.22-1.67,126.33-1.67,84.23.53,126.34,1.74S496.45,272.08,538.55,275.5Z")
                    path(fill=color, d="M338.13,221.33c-7.79,10.41-16.23,20.16-24.87,29.71s-17.58,18.8-26.73,27.84S268,296.76,258.45,305.35,239,322.29,228.6,330.05c7.84-10.37,16.29-20.11,24.92-29.67s17.57-18.81,26.68-27.88,18.44-17.94,28-26.52S327.66,229,338.13,221.33Z")
                    path(fill=color, d="M338.76,331c-10.41-7.79-20.15-16.24-29.69-24.89s-18.79-17.6-27.82-26.75-17.86-18.51-26.45-28.11-16.92-19.43-24.67-29.86c10.35,7.84,20.09,16.3,29.64,24.94s18.8,17.58,27.86,26.7,17.93,18.45,26.51,28.06S331.06,320.53,338.76,331Z")
                    path(fill="none", stroke=color, stroke-miterlimit="10", stroke-width="30", d="M138.39,272.21A143.53,143.53,0,0,1,281.91,130.89")
                    path(fill="none", stroke=color, stroke-miterlimit="10", stroke-width="10", d="M284.89,417.93c-1,0-2,0-3,0A143.53,143.53,0,0,1,138.37,274.43")
                    path(fill="none", stroke=color, stroke-miterlimit="10", stroke-width="30", d="M425.44,274.43a143.53,143.53,0,0,1-140.55,143.5")
                    path(fill="none", stroke=color, stroke-miterlimit="10", stroke-width="10", d="M284.89,130.92A143.54,143.54,0,0,1,425.43,272.21")
                }
                h2(class="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white") {
                    "Login to Star Scope"
                }
            }
            div(class="mt-8 sm:mx-auto sm:w-full sm:max-w-md") {
                div(class="bg-white dark:bg-purple-700 py-8 px-4 shadow sm:rounded-lg sm:px-10") {
                    div(class="flex items-center justify-center space-y-6") {
                        button(on:click=move |_| {
                        }, class=format!("w-[238px] border-2 border-keplr-blue-500 hover:border-keplr-blue-600 {} {}", class_button, if is_mobile { "hidden" } else { "" })) {
                            span(class="mr-2") {KeplrSvg{}}
                            "Login with Keplr"
                        }
                        p(class=format!("bg-keplr-blue-500 hover:bg-keplr-blue-600 {}", if is_mobile { "" } else {"hidden"})) { "Mobile devices are not supported yet" }
                    }
                    div(class="flex items-center justify-center space-y-6 mt-6") {
                        a(class=format!("w-[238px] bg-discord-purple-500 hover:bg-discord-purple-600 {}", class_button), href=discord_login_url) {
                            span(class="w-6 h-6 mr-2 icon-[mingcute--discord-fill]") {}
                            "Login with Discord"
                        }
                    }
                    div(class="flex items-center justify-center space-y-6 mt-6") {
                        script(async=true, src="https://telegram.org/js/telegram-widget.js?22",
                            data-telegram-login=keys::TELEGRAM_BOT_NAME,
                            data-size="large",
                            data-radius="10",
                            data-auth-url=keys::WEB_APP_URL,
                            data-request-access="write") {}
                    }
                }
            }
        }
    )
}
