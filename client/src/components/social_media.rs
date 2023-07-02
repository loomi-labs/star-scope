use log::debug;
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;
use urlencoding::encode;
use wasm_bindgen::{JsCast, UnwrapThrowExt};
use wasm_bindgen::prelude::wasm_bindgen;
use web_sys::{HtmlInputElement, MessageEvent, Window};

use crate::{AppState, AuthState, Services};
use crate::components::messages::create_error_msg_from_status;
use crate::config::keys;

#[derive(Prop)]
pub struct TelegramLoginButtonProps<'a> {
    #[builder(default = keys::WEB_APP_URL.to_string())]
    web_app_url: String,
    #[builder(default = None)]
    is_hidden: Option<&'a ReadSignal<bool>>,
}

fn hide_telegram_login_button(cx: Scope) {
    spawn_local_scoped(cx, async move {
        gloo_timers::future::TimeoutFuture::new(100).await;
        let element_id = format!("telegram-login-{}", keys::TELEGRAM_BOT_NAME);
        let element = web_sys::window()
            .unwrap()
            .document()
            .unwrap()
            .get_element_by_id(element_id.as_str());
        if element.is_some() {
            let _ = element.unwrap().set_attribute("hidden", "");
        } else {
            hide_telegram_login_button(cx)
        }
    });
}

#[component]
pub fn TelegramLoginButton<'a, G: Html>(
    cx: Scope<'a>,
    props: TelegramLoginButtonProps<'a>,
) -> View<G> {
    if let Some(is_hidden) = props.is_hidden {
        create_effect(cx, move || {
            if *is_hidden.get() {
                // necessary because the telegram script inserts the button no matter what
                hide_telegram_login_button(cx.clone());
            }
        });
    }

    view!(
        cx,
        script(
            async=true,
            src="https://telegram.org/js/telegram-widget.js?22",
            data-telegram-login=keys::TELEGRAM_BOT_NAME,
            data-size="large",
            data-radius="10",
            data-userpic="false",
            data-auth-url=props.web_app_url,
            data-request-access="write",
        ) {}
    )
}

#[derive(Prop)]
pub struct DiscordLoginButtonProps {
    text: String,
    #[builder(default = keys::WEB_APP_URL.to_string())]
    web_app_url: String,
    #[builder(default = false)]
    open_in_new_tab: bool,
}

#[component]
pub fn DiscordLoginButton<G: Html>(cx: Scope, props: DiscordLoginButtonProps) -> View<G> {
    let class_button = "flex items-center justify-center py-2 px-4 rounded-md shadow-sm text-sm \
    font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500";

    let discord_login_url = format!(
        "https://discord.com/api/oauth2/authorize?client_id={}&redirect_uri={}&response_type=code&scope=identify",
        keys::DISCORD_CLIENT_ID,
        encode(props.web_app_url.as_str())
    );

    let target = if props.open_in_new_tab {
        "_blank"
    } else {
        "_self"
    };

    view!(
        cx,
        a(class=format!("w-[219px] bg-discord-purple-500 hover:bg-discord-purple-600 {}", class_button), href=discord_login_url, target=target) {
                            span(class="w-6 h-6 mr-2 icon-[mingcute--discord-fill]") {}
                            (props.text)
        }
    )
}

const IFRAME_INPUT_ID: &str = "iframe-input";

pub fn setup_iframe_message_listener() {
    let window = web_sys::window().expect("Missing Window");
    let location = window.location();
    let host = location.hostname().unwrap();

    gloo_events::EventListener::new(&window, "message", move |event| {
        if let Some(message_event) = event.dyn_ref::<MessageEvent>() {
            if message_event.origin().contains(host.as_str()) || (message_event.origin().contains("localhost") && host == "127.0.0.1") {
                if let Some(data) = message_event.data().as_string() {
                    let window = web_sys::window().expect("Missing Window");
                    let hidden_input = window.document()
                        .and_then(|document| document.get_element_by_id(IFRAME_INPUT_ID))
                        .and_then(|input| input.dyn_into::<web_sys::HtmlInputElement>().ok()).unwrap();
                    hidden_input.set_value(data.as_str());
                }
            }
        }
    }).forget();
}

async fn login_with_wallet(cx: Scope<'_>, login_str: String) -> Result<(), ()> {
    let app_state = use_context::<AppState>(cx);
    let response = use_context::<Services>(cx)
        .auth_manager
        .clone()
        .login(login_str.clone())
        .await;
    match response {
        Ok(_) => {
            let mut auth_state = app_state.auth_state.modify();
            *auth_state = AuthState::LoggedIn;
            Ok(())
        }
        Err(status) => {
            create_error_msg_from_status(cx, status);
            Err(())
        },
    }
}

// Hack to read the value of the hidden input
fn start_login_input_timer(cx: Scope) {
    spawn_local_scoped(cx, async move {
        gloo_timers::future::TimeoutFuture::new(200).await;
        let window = web_sys::window().expect("Missing Window");
        let input_field = window.document()
            .and_then(|document| document.get_element_by_id(IFRAME_INPUT_ID))
            .and_then(|input| input.dyn_into::<HtmlInputElement>().ok()).unwrap();
        let value = input_field.value();
        if !value.is_empty() {
            spawn_local_scoped(cx.clone(), async move {
                if login_with_wallet(cx, value).await.is_err(){
                    input_field.set_value("");
                    start_login_input_timer(cx);
                }
            });
        } else {
            start_login_input_timer(cx);
        }
    });
}

#[component]
pub fn CosmosLoginButton<G: Html>(cx: Scope) -> View<G> {
    setup_iframe_message_listener();
    start_login_input_timer(cx);

    view!(
        cx,
        input(id=IFRAME_INPUT_ID, class="text-black", type="text", value="", hidden=true) {}
        iframe(
            id="iframe", 
            class="w-full h-full", src="http://localhost:3000") {}
    )
}