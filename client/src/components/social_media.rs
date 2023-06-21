use sycamore::prelude::*;
use urlencoding::encode;

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

    let target = if props.open_in_new_tab { "_blank" } else { "_self" };

    view!(
        cx,
        a(class=format!("w-[219px] bg-discord-purple-500 hover:bg-discord-purple-600 {}", class_button), href=discord_login_url, target=target) {
                            span(class="w-6 h-6 mr-2 icon-[mingcute--discord-fill]") {}
                            (props.text)
        }
    )
}

