use log::debug;
use serde::__private::de;
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;
use urlencoding::encode;
use wasm_bindgen::JsCast;
use wasm_bindgen::prelude::{wasm_bindgen, Closure};
use web_sys::{Event, MessageEvent, Window};
use web_sys::console::debug;

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

#[wasm_bindgen]
extern "C" {
    #[wasm_bindgen(js_namespace = console)]
    fn log(s: &str);
    #[wasm_bindgen(js_namespace = window)]
    fn window() -> Window;
}

#[wasm_bindgen]
pub fn setup_message_listener() {
    debug!("setting up message listener");
    let window = web_sys::window().expect("Missing Window");
    let location = window.location();
    let host = location.hostname().unwrap();
    let closure = Closure::wrap(Box::new(move |event: Event| {
        debug!("received event");
        if let Some(message_event) = event.dyn_ref::<MessageEvent>() {
            if message_event.origin().contains(host.as_str()) {
                let data = message_event.data();
                debug!("data is {:?}", data);
            }
        }
    }) as Box<dyn FnMut(_)>);

    window
        .add_event_listener_with_callback("message", closure.as_ref().unchecked_ref())
        .unwrap();

    closure.forget();
}

#[component]
pub fn CosmosLoginButton<G: Html>(cx: Scope) -> View<G> {
    setup_message_listener();

    view!(
        cx,
        iframe(
            id="iframe", 
            class="w-full h-full", src="http://localhost:3000",  on:message=move |_| {
            debug!("received message from iframe");
        }) {}
    )
}