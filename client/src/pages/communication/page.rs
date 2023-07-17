use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;

use crate::components::loading::LoadingSpinner;
use crate::components::messages::{create_error_msg_from_status, create_message};
use crate::components::social_media::{DiscordLoginButton, TelegramLoginButton};
use crate::config::keys;
use crate::types::protobuf::grpc_auth::{
    ConnectDiscordRequest, ConnectTelegramRequest,
    DiscordLoginRequest,
    TelegramLoginRequest,
};
use crate::types::protobuf::grpc_user::{DiscordChannel, DeleteDiscordChannelRequest, TelegramChat, DeleteTelegramChatRequest};
use crate::utils::url::{clean_query_params, get_discord_login_data, get_telegram_login_data};
use crate::{query_user_info, AppRoutes, AppState, InfoLevel, Services};

async fn query_discord_channels(
    cx: Scope<'_>,
    channels: &Signal<Vec<DiscordChannel>>,
    is_loading: &Signal<bool>,
) {
    is_loading.set(true);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(());
    let response = services
        .grpc_client
        .get_user_service()
        .list_discord_channels(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(result) = response {
        channels.set(result.channels);
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
    is_loading.set(false);
}

async fn delete_discord_channel(
    cx: Scope<'_>,
    channels: &Signal<Vec<DiscordChannel>>,
    channel_id: i64,
    is_loading: &Signal<bool>,
) {
    is_loading.set(true);
    let services = use_context::<Services>(cx);
    let request = services
        .grpc_client
        .create_request(DeleteDiscordChannelRequest { channel_id });
    let response = services
        .grpc_client
        .get_user_service()
        .delete_discord_channel(request)
        .await
        .map(|res| res.into_inner());
    if response.is_ok() {
        channels.set(
            channels
                .get_untracked()
                .iter()
                .filter(|channel| channel.channel_id != channel_id)
                .cloned()
                .collect(),
        );
        create_message(
            cx,
            "Channel deleted",
            "Channel successfully deleted".to_string(),
            InfoLevel::Success,
        );
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
    is_loading.set(false);
}

#[component(inline_props)]
pub fn AddEntityDialog<'a, G: Html>(
    cx: Scope<'a>,
    is_open: &'a Signal<bool>,
    service_name: &'a str,
    entity_name: &'a str,
    icon: &'a str,
    icon_bg_color: &'a str,
) -> View<G> {
    let app_state = use_context::<AppState>(cx);

    create_effect(cx, move || {
        if *is_open.get() {
            app_state.set_showing_dialog(true); // sets the backdrop to be visible
        }
    });

    create_effect(cx, move || {
        if !(*app_state.is_dialog_open.get()) {
            is_open.set(false);
        }
    });

    view! {cx,
        dialog(class="fixed inset-0 bg-white p-4 rounded-lg z-40", open=*is_open.get()) {
            div(class="flex flex-col p-4") {
                div(class="flex flex-col items-center") {
                    div(class=format!("flex items-center justify-center rounded-full text-white w-8 h-8 {}", icon_bg_color)) {
                        span(class=format!("w-6 h-6 {}", icon)) {}
                    }
                    h2(class="text-lg font-semibold mt-2") { (format!("Add {} {}", service_name, entity_name)) }
                    p(class="mt-4 text-center") { (format!("Add the {} bot to your {}.", service_name, entity_name)) }
                    p(class="mb-4 text-center") { (format!("Then send the bot command `/start`.")) }
                }
                div(class="flex justify-center mt-2") {
                    button(class="border-2 border-gray-500 text-gray-500 hover:bg-gray-500 hover:text-white font-semibold px-4 py-2 rounded mr-2",
                            on:click=move |event: web_sys::Event| {
                        event.stop_propagation();
                        app_state.set_showing_dialog(false);
                    }) { "Okay" }
                }
            }
        }
    }
}

#[component(inline_props)]
pub fn DeleteEntityDialog<'a, G: Html>(
    cx: Scope<'a>,
    is_open: &'a Signal<Option<i64>>,
    delete_signal: &'a Signal<Option<i64>>,
    service_name: &'a str,
    entity_name: &'a str,
) -> View<G> {
    let app_state = use_context::<AppState>(cx);

    create_effect(cx, move || {
        if is_open.get().is_some() {
            app_state.set_showing_dialog(true); // sets the backdrop to be visible
        }
    });

    create_effect(cx, move || {
        if !(*app_state.is_dialog_open.get()) {
            is_open.set(None);
        }
    });

    view! {cx,
        dialog(class="fixed inset-0 bg-white p-4 rounded-lg z-40", open=is_open.get().is_some()) {
            div(class="flex flex-col p-4") {
                div(class="flex flex-col items-center") {
                    span(class="w-12 h-12 text-black icon-[ph--trash]") {}
                    h2(class="text-lg font-semibold") { (format!("Delete {} {}", service_name, entity_name)) }
                    p(class="my-4 text-center") { (format!("Are you sure you want to delete this {}?", entity_name)) }
                }
                div(class="flex justify-center mt-2") {
                    button(class="border-2 border-gray-500 text-gray-500 hover:bg-gray-500 hover:text-white font-semibold px-4 py-2 rounded mr-2",
                            on:click=move |event: web_sys::Event| {
                        event.stop_propagation();
                        app_state.set_showing_dialog(false);
                    }) { "Cancel" }
                    button(class="bg-red-500 hover:bg-red-600 text-white font-semibold px-4 py-2 rounded",
                            on:click=move |event: web_sys::Event| {
                        event.stop_propagation();
                        if let Some(channel_id) = *is_open.get() {
                            delete_signal.set(Some(channel_id));
                        }
                        app_state.set_showing_dialog(false);
                    }) { "Delete" }
                }
            }
        }
    }
}

const CARD_DIV_CLASS: &str = "flex flex-col w-full";
const CARD_TITLE_CLASS: &str = "text-2xl font-semibold";
const CARD_SUBTITLE_CLASS: &str = "text-base font-semibold mt-2";
const CARD_LIST_UL_CLASS: &str = "space-y-2 mt-4";
const CARD_LIST_LI_CLASS: &str = "border-b border-gray-200 dark:border-purple-600";
const CARD_LIST_LI_ROW_CLASS: &str = "flex items-center justify-items-start my-2";
const CARD_LIST_LI_ROW_NAME_CLASS: &str = "flex-grow";
const CARD_LIST_LI_ROW_DELETE_BUTTON_CLASS: &str =
    "flex items-center ml-4 bg-red-500 hover:bg-red-600 text-white font-semibold px-1 py-1 rounded";

const CARD_ADD_DIV: &str = "flex items-center justify-items-end mt-4";
const CARD_ADD_DIV_BUTTON: &str = "flex items-center justify-center py-2 px-4 rounded-md shadow-sm text-sm font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500";

#[derive(Prop)]
pub struct DiscordCardProps {
    web_app_url: String,
    #[builder(default = false)]
    center_button: bool,
}

#[component]
pub async fn DiscordCard<G: Html>(cx: Scope<'_>, props: DiscordCardProps) -> View<G> {
    let app_state = use_context::<AppState>(cx);

    let is_connected = create_selector(cx, move || {
        app_state
            .user
            .get()
            .as_ref()
            .clone()
            .map(|user| user.has_discord)
            .unwrap_or_else(|| false)
    });

    let channels = create_signal(cx, Vec::<DiscordChannel>::new());
    let is_loading = create_signal(cx, false);
    let show_add_channel_dialog = create_signal(cx, false);
    let show_delete_dialog = create_signal(cx, None::<i64>);
    let delete_signal = create_signal(cx, None::<i64>);

    if let Some(req) = get_discord_login_data() {
        let web_app_url = props.web_app_url.clone();
        spawn_local_scoped(cx, async move {
            connect_discord_account(cx, req, web_app_url, is_loading).await;
        });
    }

    create_effect(cx, move || {
        if *is_connected.get() && !(*show_add_channel_dialog.get()) {
            // query channels if the user has discord and the add channel dialog is not open
            // -> query once when the user connects discord
            // -> queries every time the user closes the add channel dialog
            spawn_local_scoped(cx, async move {
                query_discord_channels(cx, channels, is_loading).await;
            });
        }
    });

    create_effect(cx, move || {
        if let Some(delete_id) = *delete_signal.get() {
            spawn_local_scoped(cx, async move {
                delete_signal.set(None);
                delete_discord_channel(cx, channels, delete_id, is_loading).await;
                query_discord_channels(cx, channels, is_loading).await;
            });
        }
    });

    view! {cx,
        AddEntityDialog(is_open=show_add_channel_dialog, service_name="Discord", entity_name="channel", icon="icon-[mingcute--discord-fill]", icon_bg_color="bg-discord-purple-500")
        DeleteEntityDialog(is_open=show_delete_dialog, delete_signal=delete_signal, service_name="Discord", entity_name="channel")
        div(class=format!("{} {}", CARD_DIV_CLASS, if props.center_button { "items-center" } else { "" })) {
            h2(class=CARD_TITLE_CLASS) { "Discord" }
            (if *is_loading.get() {
                view! {cx,
                    LoadingSpinner {}
                }
            } else {
                match *is_connected.get() {
                    false => {
                        let web_app_url = props.web_app_url.clone();
                        view! {cx,
                            p(class="my-4") { "Receive notifications via Discord." }
                            DiscordLoginButton(text="Connect Discord".to_string(), open_in_new_tab=false, web_app_url=web_app_url)
                        }
                    }
                    true => {
                        let discord_login_url = format!(
                            "https://discord.com/api/oauth2/authorize?client_id={}&permissions=2048&scope=bot",
                            keys::DISCORD_CLIENT_ID,
                        );
                        view! {cx,
                            (if channels.get().is_empty() {
                                view!{cx,
                                    p(class="mt-4") { "You have to add at least one channel." }
                                }
                            } else {
                                view!{cx,
                                    h3(class=CARD_SUBTITLE_CLASS) { "Connected Channels" }
                                }
                            })

                            ul(class=CARD_LIST_UL_CLASS) {
                                Indexed(
                                    iterable=channels,
                                    view=move|cx, channel| {
                                        view! { cx,
                                            li(class=CARD_LIST_LI_CLASS) {
                                                div(class=CARD_LIST_LI_ROW_CLASS) {
                                                    p(class=CARD_LIST_LI_ROW_NAME_CLASS) { (format!("{}{}", channel.name, if channel.is_group { " (Group)"} else {""} )) }
                                                    button(class=CARD_LIST_LI_ROW_DELETE_BUTTON_CLASS,
                                                        on:click=move |event: web_sys::Event| {
                                                            event.stop_propagation();
                                                            show_delete_dialog.set(Some(channel.channel_id));
                                                    }) {
                                                        span(class="w-6 h-6 icon-[ph--trash]") {}
                                                    }
                                                }
                                            }
                                        }
                                    }
                                )
                            }
                            div(class=CARD_ADD_DIV) {
                                a(class=format!("w-48 bg-discord-purple-500 hover:bg-discord-purple-600 {}", CARD_ADD_DIV_BUTTON), href=discord_login_url, target="_blank",
                                        on:click=move |event: web_sys::Event| {
                                    event.stop_propagation();
                                    show_add_channel_dialog.set(true);
                                }) {
                                    span(class="w-6 h-6 mr-2 icon-[mingcute--discord-fill]") {}
                                    "Add Channel"
                                }
                            }
                        }
                    }
                }
            })
        }
    }
}

async fn query_telegram_chats(
    cx: Scope<'_>,
    chats: &Signal<Vec<TelegramChat>>,
    is_loading: &Signal<bool>,
) {
    is_loading.set(true);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(());
    let response = services
        .grpc_client
        .get_user_service()
        .list_telegram_chats(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(result) = response {
        chats.set(result.chats);
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
    is_loading.set(false);
}

async fn delete_telegram_chat(
    cx: Scope<'_>,
    chats: &Signal<Vec<TelegramChat>>,
    chat_id: i64,
    is_loading: &Signal<bool>,
) {
    is_loading.set(true);
    let services = use_context::<Services>(cx);
    let request = services
        .grpc_client
        .create_request(DeleteTelegramChatRequest { chat_id });
    let response = services
        .grpc_client
        .get_user_service()
        .delete_telegram_chat(request)
        .await
        .map(|res| res.into_inner());
    if response.is_ok() {
        chats.set(
            chats
                .get_untracked()
                .iter()
                .filter(|chat| chat.chat_id != chat_id)
                .cloned()
                .collect(),
        );
        create_message(
            cx,
            "Chat deleted",
            "Chat successfully deleted".to_string(),
            InfoLevel::Success,
        );
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
    is_loading.set(false);
}

#[component(inline_props)]
pub async fn TelegramCard<G: Html>(
    cx: Scope<'_>,
    web_app_url: String,
    center_button: bool,
) -> View<G> {
    let app_state = use_context::<AppState>(cx);

    let is_connected = create_selector(cx, move || {
        app_state
            .user
            .get()
            .as_ref()
            .clone()
            .map(|user| user.has_telegram)
            .unwrap_or_else(|| false)
    });

    let chats = create_signal(cx, Vec::<TelegramChat>::new());
    let is_loading = create_signal(cx, false);
    let show_delete_dialog = create_signal(cx, None::<i64>);
    let show_add_chat_dialog = create_signal(cx, false);
    let delete_signal = create_signal(cx, None::<i64>);

    if let Some(req) = get_telegram_login_data() {
        spawn_local_scoped(cx, async move {
            connect_telegram_account(cx, req, is_loading).await;
        });
    }

    create_effect(cx, move || {
        if *is_connected.get() && !(*show_add_chat_dialog.get()) {
            // query chats if the user has telegram and the add chat dialog is not open
            // -> query once the user has connected telegram
            // -> queries every time the add chat dialog is closed
            spawn_local_scoped(cx, async move {
                query_telegram_chats(cx, chats, is_loading).await;
            });
        }
    });

    create_effect(cx, move || {
        if let Some(delete_id) = *delete_signal.get() {
            spawn_local_scoped(cx, async move {
                delete_signal.set(None);
                delete_telegram_chat(cx, chats, delete_id, is_loading).await;
                query_telegram_chats(cx, chats, is_loading).await;
            });
        }
    });

    view! {cx,
        AddEntityDialog(is_open=show_add_chat_dialog, service_name="Telegram", entity_name="chat", icon="icon-[bxl--telegram]", icon_bg_color="bg-telegram-blue-500")
        DeleteEntityDialog(is_open=show_delete_dialog, delete_signal=delete_signal, service_name="Telegram", entity_name="chat")
        div(class=format!("{} {}", CARD_DIV_CLASS, if center_button { "items-center" } else { "" })) {
            h2(class="text-2xl font-semibold") { "Telegram" }
            (if *is_loading.get() {
                view! {cx,
                    LoadingSpinner {}
                }
            } else {
                match *is_connected.get() {
                    false => {
                        let web_app_url = web_app_url.clone();
                        view! {cx,
                            p(class="my-4") { "Receive notifications via Telegram." }
                            TelegramLoginButton(web_app_url=web_app_url, is_hidden=Some(is_connected))
                        }
                    }
                    true => {
                        let tg_bot_url = format!("https://t.me/{}", keys::TELEGRAM_BOT_NAME);
                        view! {cx,
                            (if chats.get().is_empty() {
                                view!{cx,
                                    p(class="mt-4") { "You have to add at least one chat." }
                                }
                            } else {
                                view!{cx,
                                    h3(class=CARD_SUBTITLE_CLASS) { "Connected Chats" }
                                }
                            })
                            ul(class=CARD_LIST_UL_CLASS) {
                                Indexed(
                                    iterable=chats,
                                    view=move|cx, chat| {
                                        view! { cx,
                                            li(class=CARD_LIST_LI_CLASS) {
                                                div(class=CARD_LIST_LI_ROW_CLASS) {
                                                    p(class=CARD_LIST_LI_ROW_NAME_CLASS) { (format!("{}{}", chat.name, if chat.is_group { " (Group)"} else {""} )) }
                                                    button(class=CARD_LIST_LI_ROW_DELETE_BUTTON_CLASS, on:click=move |_| show_delete_dialog.set(Some(chat.chat_id))) {
                                                        span(class="w-6 h-6 icon-[ph--trash]") {}
                                                    }
                                                }
                                            }
                                        }
                                    }
                                )
                            }
                            div(class=CARD_ADD_DIV) {
                                a(class=format!("w-48 bg-telegram-blue-500 hover:bg-telegram-blue-600 {}", CARD_ADD_DIV_BUTTON), href=tg_bot_url, target="_blank", on:click=move |_| show_add_chat_dialog.set(true)) {
                                    span(class="w-6 h-6 mr-2 icon-[bxl--telegram]") {}
                                    "Add Chat"
                                }
                            }
                        }
                    }
                }
            })
        }
    }
}

async fn connect_discord_account(
    cx: Scope<'_>,
    data: DiscordLoginRequest,
    web_app_url: String,
    is_loading: &Signal<bool>,
) {
    is_loading.set(true);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(ConnectDiscordRequest {
        code: data.code,
        web_app_url,
    });
    let response = services
        .grpc_client
        .get_auth_service()
        .connect_discord(request)
        .await
        .map(|res| res.into_inner());
    if response.is_ok() {
        query_user_info(cx).await;
        create_message(
            cx,
            "Discord account connected",
            "Your Discord account is now connected. Add channels to receive notifications.",
            InfoLevel::Success,
        );
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
    clean_query_params();
    is_loading.set(false);
}

async fn connect_telegram_account(
    cx: Scope<'_>,
    data: TelegramLoginRequest,
    is_loading: &Signal<bool>,
) {
    is_loading.set(true);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(ConnectTelegramRequest {
        data_str: data.data_str,
        hash: data.hash,
    });
    let response = services
        .grpc_client
        .get_auth_service()
        .connect_telegram(request)
        .await
        .map(|res| res.into_inner());
    if response.is_ok() {
        query_user_info(cx).await;
        create_message(
            cx,
            "Telegram account connected",
            "Your Telegram account is now connected. Add chats to receive notifications.",
            InfoLevel::Success,
        );
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
    clean_query_params();
    is_loading.set(false);
}

#[component]
pub fn Communication<G: Html>(cx: Scope) -> View<G> {
    let web_app_url = create_ref(
        cx,
        keys::WEB_APP_URL.to_string() + AppRoutes::Communication.to_string().as_str(),
    );

    view! {cx,
        div(class="container mx-auto") {
            div(class="flex flex-col space-y-4") {
                div {
                    h1(class="text-4xl font-bold") { "Communication Channels" }
                }
                div(class="p-8 rounded-lg dark:bg-purple-700") {
                    DiscordCard(web_app_url=web_app_url.clone())
                }
                div(class="p-8 rounded-lg dark:bg-purple-700") {
                    TelegramCard(web_app_url=web_app_url.clone(), center_button=false)
                }
            }
        }
    }
}
