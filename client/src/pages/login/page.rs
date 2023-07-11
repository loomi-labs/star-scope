use sycamore::{prelude::*, view};
use wasm_bindgen::prelude::*;

use crate::components::social_media::{CosmosLoginButton, DiscordLoginButton, TelegramLoginButton};
use crate::config::keys;

#[wasm_bindgen(module = "/src/pages/login/wallet_connect_login.js")]
extern "C" {
    fn isMobile(url: String) -> JsValue;
}

fn is_mobile() -> bool {
    let result = isMobile(keys::WEB_APP_URL.to_string());
    serde_wasm_bindgen::from_value(result).unwrap_or(false)
}

#[component]
pub async fn Login<G: Html>(cx: Scope<'_>) -> View<G> {
    let color = keys::WHITE_COLOR;

    let is_mobile = is_mobile();

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
            div(class="mt-8 mx-auto w-full h-full max-w-md") {
                div(class="py-8 px-4 shadow rounded-lg px-10") {
                    div(class="flex items-center justify-center space-y-6 mt-6") {
                        DiscordLoginButton(text="Login with Discord".to_string())
                    }
                    div(class="flex items-center justify-center space-y-6 mt-6") {
                        TelegramLoginButton(web_app_url=keys::WEB_APP_URL.to_string())
                    }
                    (if !is_mobile {
                        view!(cx,
                            div(class="flex justify-center space-y-6 mt-6") {
                                CosmosLoginButton()
                            }
                        )
                    } else {
                        view!(cx,
                            div()
                        )
                    })
                }
            }
        }
    )
}
