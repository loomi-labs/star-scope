use sycamore::{prelude::*, view};
use wasm_bindgen::prelude::*;

use crate::components::images::StarScopeLogo;
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
    let is_mobile = is_mobile();

    view!(cx,
        div(class="h-screen w-screen bg-gray-100 dark:bg-d-bg flex flex-col justify-center py-12 sm:px-6 lg:px-8") {
            div(class="sm:mx-auto sm:w-full sm:max-w-md flex flex-col justify-center items-center") {
                StarScopeLogo(width=200, height=200, color=keys::WHITE_COLOR)
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
