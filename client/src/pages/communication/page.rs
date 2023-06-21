use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;

use crate::{AppRoutes, AppState, InfoLevel, query_user_info, Services};
use crate::components::loading::LoadingSpinner;
use crate::components::messages::{create_error_msg_from_status, create_message};
use crate::components::social_media::{DiscordLoginButton, TelegramLoginButton};
use crate::config::keys;
use crate::types::protobuf::grpc::{ConnectDiscordRequest, DeleteDiscordChannelRequest, DiscordChannel};
use crate::utils::url::get_query_param;


async fn query_discord_channels(cx: Scope<'_>, channels: &Signal<Vec<DiscordChannel>>, is_loading: &Signal<bool>) {
    is_loading.set(true);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(());
    let response = services
        .grpc_client
        .get_user_service()
        .get_discord_channels(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(result) = response {
        channels.set(result.channels);
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
    is_loading.set(false);
}

async fn delete_discord_channel(cx: Scope<'_>, channels: &Signal<Vec<DiscordChannel>>, channel_id: i64, is_loading: &Signal<bool>) {
    is_loading.set(true);
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(DeleteDiscordChannelRequest { channel_id });
    let response = services
        .grpc_client
        .get_user_service()
        .delete_discord_channel(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(_) = response {
        channels.set(channels.get_untracked().iter().filter(|channel| channel.channel_id != channel_id).cloned().collect());
        create_message(cx, "Channel deleted", "Channel successfully deleted".to_string(), InfoLevel::Success);
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
    is_loading.set(false);
}

#[component]
pub async fn DiscordCard<G: Html>(cx: Scope<'_>) -> View<G> {
    let app_state = use_context::<AppState>(cx);

    let is_connected = create_selector(cx, move || {
        app_state.user.get().as_ref().clone().map(|user| user.has_discord).unwrap_or_else(|| false)
    });

    let channels = create_signal(cx, Vec::<DiscordChannel>::new());
    let is_loading = create_signal(cx, false);
    let show_delete_dialog = create_signal(cx, None::<i64>);
    let show_add_channel_dialog = create_signal(cx, false);

    create_effect(cx, move || {
        if *is_connected.get() {
            spawn_local_scoped(cx, async move {
                query_discord_channels(cx.clone(), &channels, &is_loading).await;
            });
        }
    });

    let class_button = "flex items-center justify-center py-2 px-4 rounded-md shadow-sm text-sm \
    font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500";

    view! {cx,
        div(class="p-8 rounded-lg dark:bg-purple-700") {
            dialog(class="absolute rounded-lg drop-shadow-lg dark:bg-white", open=*show_add_channel_dialog.get()) {
                div(class="flex flex-col p-4") {
                    div(class="flex flex-col items-center") {
                        div(class="flex items-center justify-center rounded-full bg-discord-purple-500 text-white w-8 h-8 mr-2") {
                            span(class="w-6 h-6 icon-[mingcute--discord-fill]") {}
                        }
                        h2(class="text-lg font-semibold mt-2") { "Add Channel" }
                        p(class="mt-4 text-center") { "Add the Star Scope bot to your server." }
                        p(class="mb-4 text-center") { "Then send the bot command `/start` in the channel(s) that you want to add." }
                    }
                    div(class="flex justify-center mt-2") {
                        button(class="border-2 border-gray-500 text-gray-500 hover:bg-gray-500 hover:text-white font-semibold px-4 py-2 rounded mr-2", on:click=move |_| {
                            spawn_local_scoped(cx, async move {
                                query_discord_channels(cx.clone(), channels, &is_loading).await;
                            });
                            show_add_channel_dialog.set(false);
                        }) { "Okay" }
                    }
                }
            }
            dialog(class="absolute rounded-lg drop-shadow-lg dark:bg-white", open=show_delete_dialog.get().is_some()) {
                div(class="flex flex-col p-4") {
                    div(class="flex flex-col items-center") {
                        span(class="w-12 h-12 text-black icon-[ph--trash]") {}
                        h2(class="text-lg font-semibold") { "Delete Channel" }
                        p(class="my-4 text-center") { "Are you sure you want to delete this channel?" }
                    }
                    div(class="flex justify-center") {
                        button(class="border-2 border-gray-500 text-gray-500 hover:bg-gray-500 hover:text-white font-semibold px-4 py-2 rounded mr-2", on:click=move |_| show_delete_dialog.set(None)) { "Cancel" }
                        button(class="bg-red-500 hover:bg-red-600 text-white font-semibold px-4 py-2 rounded", on:click=move |_| {
                            let channel_id = show_delete_dialog.get().unwrap();
                            spawn_local_scoped(cx, async move {
                                delete_discord_channel(cx.clone(), channels, channel_id, &is_loading).await;
                            });
                            show_delete_dialog.set(None);
                        }) { "Delete" }
                    }
                }
            }
            div {
                h2(class="text-2xl font-semibold") { "Discord" }
                (match *is_connected.get() {
                    false => {
                        let web_app_url = keys::WEB_APP_URL.to_string() + AppRoutes::Communication.to_string().as_str();
                        view! {cx,
                            p(class="my-4") { "Connect your Discord account to receive notifications via Discord." }
                            DiscordLoginButton(text="Connect Discord".to_string(), open_in_new_tab=false, web_app_url=web_app_url)
                        }
                    }
                    true => {
                        let discord_login_url = format!(
                            "https://discord.com/api/oauth2/authorize?client_id={}&permissions=2048&scope=bot",
                            keys::DISCORD_CLIENT_ID,
                        );
                        if *is_loading.get() {
                            view! {cx,
                                LoadingSpinner {}
                            }
                        } else {
                            view! {cx,
                                h3(class="text-lg font-semibold mt-2") { "Connected Channels" }
                                ul(class="space-y-2") {
                                    Indexed(
                                        iterable=channels,
                                        view=move|cx, channel| {
                                            view! { cx,
                                                li(class="border-b border-gray-200 dark:border-purple-600") {
                                                    div(class="flex items-center justify-items-start my-2") {
                                                        p(class="flex-grow") { (format!("{}{}", channel.name, if channel.is_group { " (Group)"} else {""} )) }
                                                        button(class="flex items-center ml-4 bg-red-500 hover:bg-red-600 text-white font-semibold px-1 py-1 rounded", on:click=move |_| show_delete_dialog.set(Some(channel.channel_id))) {
                                                            span(class="w-6 h-6 icon-[ph--trash]") {}
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    )
                                }
                                div(class="flex items-center justify-items-end mt-4") {
                                    a(class=format!("w-48 bg-discord-purple-500 hover:bg-discord-purple-600 {}", class_button), href=discord_login_url, target="_blank", on:click=move |_| show_add_channel_dialog.set(true)) {
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
}

async fn connect_discord_account(cx: Scope<'_>, code: String) {
    let services = use_context::<Services>(cx);
    let web_app_url = keys::WEB_APP_URL.to_string() + AppRoutes::Communication.to_string().as_str();
    let request = services.grpc_client.create_request(ConnectDiscordRequest { code, web_app_url });
    let response = services
        .grpc_client
        .get_auth_service()
        .connect_discord(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(_) = response {
        query_user_info(cx).await;
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

#[component]
pub fn Communication<G: Html>(cx: Scope) -> View<G> {
    if let Some(code) = get_query_param("code") {
        spawn_local_scoped(cx, async move {
            connect_discord_account(cx.clone(), code).await;
        });
    }

    view! {cx,
        div(class="container mx-auto") {
            div(class="flex flex-col space-y-4") {
                div {
                    h1(class="text-4xl font-bold") { "Communication Channels" }
                }
                div {
                    // TelegramCard {}
                }
                div {
                    DiscordCard {}
                }
                // div {
                //     Card(state=CardState::Connected)
                // }
            }
        }
    }
}
