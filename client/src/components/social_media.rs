use sycamore::prelude::*;
use crate::config::keys;

#[component]
pub fn TelegramLoginButton<G: Html>(cx: Scope) -> View<G> {
    view!(
        cx,
        script(async=true, src="https://telegram.org/js/telegram-widget.js?22",
            data-telegram-login=keys::TELEGRAM_BOT_NAME,
            data-size="large",
            data-radius="10",
            data-userpic="false",
            data-auth-url=keys::WEB_APP_URL,
            data-request-access="write") {}
    )
}
