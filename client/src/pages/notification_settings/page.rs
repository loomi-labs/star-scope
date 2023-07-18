use crate::components::messages::{
    create_error_msg_from_status, create_message, create_timed_message,
};
use crate::types::protobuf::grpc_settings::{
    RemoveWalletRequest, UpdateWalletRequest, ValidateWalletRequest, Wallet,
};
use crate::{AppState, InfoLevel, Services};
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;

async fn query_wallets(cx: Scope<'_>) -> Result<Vec<Wallet>, ()> {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(());
    let response = services
        .grpc_client
        .get_settings_service()
        .get_wallets(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(result) = response {
        Ok(result.wallets)
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        Err(())
    }
}
enum Update {
    Funding,
    Staking,
    GovVotingReminder,
}

async fn update_existing_wallet(cx: Scope<'_>, wallet_sig: &Signal<Wallet>, update: Update) {
    let wallet = create_ref(cx, wallet_sig.get_untracked());
    let notify_funding = if let Update::Funding = update {
        !wallet.notify_funding
    } else {
        wallet.notify_funding
    };
    let notify_staking = if let Update::Staking = update {
        !wallet.notify_staking
    } else {
        wallet.notify_staking
    };
    let notify_gov_voting_reminder = if let Update::GovVotingReminder = update {
        !wallet.notify_gov_voting_reminder
    } else {
        wallet.notify_gov_voting_reminder
    };
    let request = UpdateWalletRequest {
        wallet_address: wallet.address.clone(),
        notify_funding,
        notify_staking,
        notify_gov_voting_reminder,
    };

    let result = update_wallet(cx, wallet_sig, request).await;
    if result.is_ok() {
        let msg = match update {
            Update::Funding => {
                if notify_funding {
                    "You will be notified about funding events"
                } else {
                    "You will no longer be notified about funding events"
                }
            }
            Update::Staking => {
                if notify_staking {
                    "You will be notified about staking events"
                } else {
                    "You will no longer be notified about staking events"
                }
            }
            Update::GovVotingReminder => {
                if notify_gov_voting_reminder {
                    "You will be notified to vote on governance proposals"
                } else {
                    "You will no longer be notified to vote on governance proposals"
                }
            }
        };
        create_timed_message(cx, "Wallet updated", msg, InfoLevel::Success, 5);
    }
}

async fn update_wallet(
    cx: Scope<'_>,
    wallet: &Signal<Wallet>,
    update: UpdateWalletRequest,
) -> Result<(), ()> {
    let update_ref = create_ref(cx, update.clone());
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(update);
    let response = services
        .grpc_client
        .get_settings_service()
        .update_wallet(request)
        .await
        .map(|res| res.into_inner());

    if response.is_ok() {
        wallet.modify().notify_funding = update_ref.notify_funding;
        wallet.modify().notify_staking = update_ref.notify_staking;
        wallet.modify().notify_gov_voting_reminder = update_ref.notify_gov_voting_reminder;
        Ok(())
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        Err(())
    }
}

const BUTTON_ROW_CLASS: &str =
    "flex items-center cursor-pointer py-1 px-2 space-x-2 rounded-lg hover:bg-purple-700";

#[component(inline_props)]
fn Tooltip<G: Html>(cx: Scope, title: &'static str) -> View<G> {
    let is_visible = create_signal(cx, false);

    let unavailable_class = "w-4 h-4 bg-red-500 icon-[iconamoon--unavailable]";

    view! {cx,
        div(class="relative") {
            div(class=BUTTON_ROW_CLASS,
                    on:mouseover=move |_| is_visible.set(true),
                    on:mouseout=move |_| is_visible.set(false)) {
                span(class=unavailable_class) {}
                span() { (title) }
                div(class=format!("absolute inset-0 italic text-xs rounded-lg dark:bg-purple-700 {}", if *is_visible.get() { "visible" } else { "invisible" })) {
                    "Not yet supported"
                }
            }
        }
    }
}

#[derive(Debug, Clone, PartialEq)]
enum WalletValidation {
    Valid(Wallet),
    Invalid(String),
}

async fn query_validate_wallet(cx: Scope<'_>, address: String) -> WalletValidation {
    let services = use_context::<Services>(cx);
    let request = services
        .grpc_client
        .create_request(ValidateWalletRequest { address });
    let response = services
        .grpc_client
        .get_settings_service()
        .validate_wallet(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(response) = response {
        if response.is_valid {
            if response.is_supported {
                if response.is_already_added {
                    return WalletValidation::Invalid("Wallet already added".to_string());
                }
                if let Some(wallet) = response.wallet {
                    return WalletValidation::Valid(wallet);
                }
                create_message(cx, "Error", "Wallet not found", InfoLevel::Error);
                WalletValidation::Invalid("Wallet not found".to_string())
            } else {
                WalletValidation::Invalid("Chain is currently not supported".to_string())
            }
        } else {
            WalletValidation::Invalid("Invalid wallet address".to_string())
        }
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        WalletValidation::Invalid("Error".to_string())
    }
}

#[component(inline_props)]
fn AddWallet<'a, G: Html>(cx: Scope<'a>, wallets: &'a Signal<Vec<&'a Signal<Wallet>>>) -> View<G> {
    let new_wallet_address = create_signal(cx, String::new());

    let validation = create_signal(cx, None::<WalletValidation>);

    let has_new_wallet = create_signal(cx, false);

    create_effect(cx, move || {
        let address = new_wallet_address.get().as_ref().clone();
        if address.is_empty() {
            validation.set(None);
        } else if address.len() < 30 {
            validation.set(Some(WalletValidation::Invalid(
                "Wallet address is too short".to_string(),
            )));
        } else if wallets.get().iter().any(|w| w.get().address == address) {
            validation.set(Some(WalletValidation::Invalid(
                "Wallet already added".to_string(),
            )));
        } else {
            spawn_local_scoped(cx, async move {
                let result =
                    query_validate_wallet(cx, new_wallet_address.get().as_ref().clone()).await;
                validation.set(Some(result.clone()));
                if let WalletValidation::Valid(wallet) = result {
                    let request = UpdateWalletRequest {
                        wallet_address: wallet.address.clone(),
                        notify_funding: true,
                        notify_staking: true,
                        notify_gov_voting_reminder: true,
                    };
                    let wallet_sig = create_signal(cx, wallet.clone());
                    if update_wallet(cx, wallet_sig, request).await.is_ok() {
                        wallets.modify().push(wallet_sig);
                        new_wallet_address.set(String::new());
                        has_new_wallet.set(false);
                        create_message(
                            cx,
                            "Wallet added",
                            format!("You added the wallet {}", wallet.address),
                            InfoLevel::Success,
                        );
                    }
                }
            });
        }
    });

    view! {cx,
        div(class=format!("flex rounded-lg p-4 mt-2 {}", if *has_new_wallet.get() { "transition ease-in-out duration-500 dark:bg-purple-800" } else { "" })) {
            div(class="flex items-center px-1 gap-1") {
                button(class=format!("flex justify-center items-center w-10 h-10 md:w-14 md:h-14 opacity-100 rounded-full {}",
                        if *has_new_wallet.get() { "bg-red-500 hover:bg-red-600" } else { "bg-green-500 hover:bg-green-600" }),
                        on:click=move |_| has_new_wallet.set(!*has_new_wallet.get())) {
                    span(class=format!("w-6 h-6 md:w-10 md:h-10 cursor-pointer {}",
                        if *has_new_wallet.get() { "icon-[iconoir--cancel]" } else { "icon-[ic--round-add]" })) {}
                }
            }
            div(class=format!("flex flex-col w-full px-2 space-y-2 {}", if *has_new_wallet.get() { "visible" } else { "invisible" })) {
                span(class="text-base font-semibold") { "Add new wallet" }
                div(class="flex flex-wrap items-center gap-x-4") {
                    div(class="flex flex-col") {
                        div(class="flex") {
                            input(
                                class="w-full placeholder:italic border border-gray-300 rounded-lg px-4 py-2 text-black focus:outline-none focus:ring-2 focus:ring-primary",
                                placeholder="Wallet address",
                                type="text",
                                bind:value=new_wallet_address,
                            )
                        }
                        (if let Some(WalletValidation::Invalid(msg)) = validation.get().as_ref().clone() {
                            view! {cx,
                                span(class="text-red-500 text-left") {(msg)}
                            }
                        } else {
                            view! {cx, }
                        })
                    }
                }
            }
        }
    }
}

async fn delete_wallet(cx: Scope<'_>, wallet_address: String) -> Result<(), ()> {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(RemoveWalletRequest {
        wallet_address: wallet_address.clone(),
    });
    let response = services
        .grpc_client
        .get_settings_service()
        .remove_wallet(request)
        .await
        .map(|res| res.into_inner());

    if response.is_ok() {
        create_timed_message(
            cx,
            "Wallet deleted",
            format!("Wallet {} was deleted", wallet_address),
            InfoLevel::Success,
            5,
        );
        Ok(())
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        Err(())
    }
}

#[component(inline_props)]
fn AskDeleteDialog<'a, G: Html>(
    cx: Scope<'a>,
    is_open: &'a Signal<Option<Wallet>>,
    delete_signal: &'a Signal<Option<String>>,
) -> View<G> {
    let app_state = use_context::<AppState>(cx);

    let handle_delete = move |wallet: Option<Wallet>| {
        is_open.set(None);
        app_state.set_showing_dialog(false);
        if wallet.is_some() {
            let wallet_address = wallet.unwrap().address;
            delete_signal.set(Some(wallet_address));
        } else {
            create_message(cx, "Error", "Wallet could not be deleted", InfoLevel::Error);
        }
    };

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
                    h2(class="text-lg font-semibold") { ("Delete wallet") }
                    p(class="my-4 text-center") { (format!("Are you sure you want to delete {}?", is_open.get().as_ref().clone().map_or("".to_string(), |w: Wallet| w.address))) }
                }
                div(class="flex justify-center mt-2") {
                    button(class="border-2 border-gray-500 text-gray-500 hover:bg-gray-500 hover:text-white font-semibold px-4 py-2 rounded mr-2",
                            on:click=move |event: web_sys::Event| {
                        event.stop_propagation();
                        app_state.set_showing_dialog(false);
                    }) { "Cancel" }
                    button(class="bg-red-500 hover:bg-red-600 text-white font-semibold px-4 py-2 rounded",
                            on:click=move |_| handle_delete(is_open.get().as_ref().clone())) { "Delete" }
                }
            }
        }
    }
}

#[component(inline_props)]
fn WalletList<'a, G: Html>(cx: Scope<'a>, wallets: &'a Signal<Vec<&'a Signal<Wallet>>>) -> View<G> {
    let selected_class = "w-4 h-4 bg-primary icon-[icon-park-solid--check-one]";
    let unselected_class = "w-4 h-4 rounded-full border-2 border-primary";

    let show_delete_dialog = create_signal(cx, None::<Wallet>);

    let delete_signal = create_signal(cx, None::<String>);

    create_effect(cx, move || {
        if let Some(wallet_address) = delete_signal.get().as_ref().clone() {
            let mut wallets: Modify<'_, Vec<&Signal<Wallet>>> = wallets.modify();
            spawn_local_scoped(cx, async move {
                if delete_wallet(cx, wallet_address.clone()).await.is_ok() {
                    wallets.retain(|w| w.get().address != wallet_address);
                }
            });
        }
    });

    view! {cx,
        div(class="flex flex-col w-full mt-2 space-y-2") {
            Indexed(
                iterable = wallets,
                view = move |cx, wallet| {
                    let cloned = wallet.get().as_ref().clone();
                    let handle_update = move |update: Update| {
                        spawn_local_scoped(cx, async move {
                            update_existing_wallet(cx, wallet, update).await;
                        });
                    };

                    let prefix = cloned.address[..8].to_owned();
                    let suffix = cloned.address[cloned.address.len() - 4..].to_owned();
                    let shortened_address = format!("{}...{}", prefix, suffix);

                    view!{cx,
                        div(class="flex p-4 rounded-lg items-center bg-purple-800") {
                            div(class="flex items-center px-1 gap-1") {
                                img(src=cloned.logo_url, class="w-10 h-10 md:w-14 md:h-14") {}
                            }
                            div(class="flex flex-col text-sm max-w-[calc(100%-theme(space.16))]") {
                                span(class="text-base font-semibold px-2") { (shortened_address) }
                                div(class="flex flex-wrap flex-shrink items-center gap-x-4") {
                                    (if wallet.get().is_notify_funding_supported {
                                        view!{cx,
                                            div(class=BUTTON_ROW_CLASS, on:click=move |_| handle_update(Update::Funding)) {
                                                span(class=(if wallet.get().notify_funding { selected_class } else { unselected_class })) {}
                                                span() { "Funding" }
                                            }
                                        }
                                    } else {
                                        view!{cx,
                                            Tooltip(title="Funding")
                                        }
                                    })
                                    (if wallet.get().is_notify_staking_supported {
                                        view!{cx,
                                            div(class=BUTTON_ROW_CLASS, on:click=move |_| handle_update(Update::Staking)) {
                                                span(class=(if wallet.get().notify_staking { selected_class } else { unselected_class })) {}
                                                span() { "Staking" }
                                            }
                                        }
                                    } else {
                                        view!{cx,
                                            Tooltip(title="Staking")
                                        }
                                    })
                                    (if wallet.get().is_notify_gov_voting_reminder_supported {
                                        view!{cx,
                                            div(class=BUTTON_ROW_CLASS, on:click=move |_| handle_update(Update::GovVotingReminder)) {
                                                span(class=(if wallet.get().notify_gov_voting_reminder { selected_class } else { unselected_class })) {}
                                                span(class="truncate") { "Governance reminders" }
                                            }
                                        }
                                    } else {
                                        view!{cx,
                                            Tooltip(title="Governance")
                                        }
                                    })
                                }
                            }
                            AskDeleteDialog(is_open=show_delete_dialog.clone(), delete_signal=delete_signal)
                            button(class="flex items-center justify-center w-10 h-10 rounded-lg dark:bg-purple-700 dark:hover:bg-purple-600",
                                on:click=move |event: web_sys::Event| {
                                    event.stop_propagation();
                                    show_delete_dialog.set(Some(wallet.get().as_ref().clone()));
                            }) {
                                span(class="w-6 h-6 icon-[ph--trash]") {}
                            }
                        }
                    }
                },
            )
        }
    }
}

#[component]
pub async fn NotificationSettings<'a, G: Html>(cx: Scope<'a>) -> View<G> {
    let wallets: &'a Signal<Vec<&'a Signal<Wallet>>> = create_signal(cx, vec![]);

    spawn_local_scoped(cx, async move {
        let result = query_wallets(cx).await;
        if let Ok(result_wallets) = result {
            let new_wallets: Vec<&Signal<Wallet>> = result_wallets
                .iter()
                .map(|wallet| create_signal(cx, wallet.clone()))
                .collect();
            wallets.set(new_wallets);
        }
    });

    view! {cx,
        div(class="flex flex-col") {
            h1(class="text-4xl font-semibold") { "Notification settings" }
            div(class="flex flex-col mt-4") {
                h2(class="text-xl font-semibold") { (format!("Wallets ({})", wallets.get().len())) }
                div(class="flex flex-col rounded-lg") {
                    WalletList(wallets=wallets)
                    AddWallet(wallets=wallets)
                }
            }
        }
    }
}
