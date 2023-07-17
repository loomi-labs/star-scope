use crate::components::messages::{
    create_error_msg_from_status, create_message, create_timed_message,
};
use crate::types::protobuf::grpc_settings::{UpdateWalletRequest, Wallet};
use crate::{InfoLevel, Services};
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
    NotifyFunding,
    NotifyStaking,
    NotifyGovVotingReminder,
}

async fn update_wallet(cx: Scope<'_>, wallet_sig: &Signal<Wallet>, update: Update) {
    let wallet = create_ref(cx, wallet_sig.get_untracked());
    let services = use_context::<Services>(cx);
    let notify_funding = if let Update::NotifyFunding = update {
        !wallet.notify_funding.clone()
    } else {
        wallet.notify_funding.clone()
    };
    let notify_staking = if let Update::NotifyStaking = update {
        !wallet.notify_staking.clone()
    } else {
        wallet.notify_staking.clone()
    };
    let notify_gov_voting_reminder = if let Update::NotifyGovVotingReminder = update {
        !wallet.notify_gov_voting_reminder.clone()
    } else {
        wallet.notify_gov_voting_reminder.clone()
    };
    let request = services.grpc_client.create_request(UpdateWalletRequest {
        wallet_address: wallet.address.clone(),
        notify_funding,
        notify_staking,
        notify_gov_voting_reminder,
    });

    let response = services
        .grpc_client
        .get_settings_service()
        .update_wallet(request)
        .await
        .map(|res| res.into_inner());

    if let Ok(_) = response {
        wallet_sig.set(Wallet {
            address: wallet.address.clone(),
            logo_url: wallet.logo_url.clone(),
            notify_funding,
            notify_staking,
            notify_gov_voting_reminder,
        });
        let msg = match update {
            Update::NotifyFunding => {
                if notify_funding {
                    "You will be notified about funding events"
                } else {
                    "You will no longer be notified about funding events"
                }
            }
            Update::NotifyStaking => {
                if notify_staking {
                    "You will be notified about staking events"
                } else {
                    "You will no longer be notified about staking events"
                }
            }
            Update::NotifyGovVotingReminder => {
                if notify_gov_voting_reminder {
                    "You will be notified to vote on governance proposals"
                } else {
                    "You will no longer be notified to vote on governance proposals"
                }
            }
        };
        create_timed_message(cx, "Wallet updated", msg, InfoLevel::Success, 5);
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

pub struct NotificationSettingsState<'a> {
    pub wallets: &'a Signal<Vec<&'a Signal<Wallet>>>,
}

#[component]
pub async fn NotificationSettings<G: Html>(cx: Scope<'_>) -> View<G> {
    let notification_settings_state = NotificationSettingsState {
        wallets: &create_signal(cx, vec![]),
    };

    spawn_local_scoped(cx, async move {
        let result = query_wallets(cx).await;
        if let Ok(wallets) = result {
            let new_wallets: Vec<&Signal<Wallet>> = wallets
                .iter()
                .map(|wallet| create_signal(cx, wallet.clone()))
                .collect();
            notification_settings_state.wallets.set(new_wallets);
        }
    });

    let selected_class = "w-4 h-4 bg-primary icon-[icon-park-solid--check-one]";
    let unselected_class = "w-4 h-4 rounded-full border-2 border-primary";
    let button_row_class =
        "flex items-center cursor-pointer p-2 space-x-2 rounded-lg hover:bg-purple-700";

    view! {cx,
        div(class="flex flex-col") {
            h1(class="text-2xl font-semibold") { "Notification settings" }
            div(class="flex flex-col mt-4") {
                h2(class="text-lg font-semibold") { (format!("Wallets ({})", notification_settings_state.wallets.get().len())) }
                div(class="flex flex-col mt-2 space-y-2") {
                    Indexed(
                        iterable = notification_settings_state.wallets,
                        view = move |cx, wallet| {
                            let cloned = wallet.get().as_ref().clone();
                            let handle_update = move |update: Update| {
                                spawn_local_scoped(cx, async move {
                                    update_wallet(cx, wallet, update).await;
                                });
                            };

                            let prefix = cloned.address[..8].to_owned();
                            let suffix = cloned.address[cloned.address.len() - 4..].to_owned();
                            let shortened_address = format!("{}...{}", prefix, suffix);

                            view!{cx,
                                div(class="flex flex-col p-4 space-y-4 rounded-lg bg-purple-800") {
                                    div(class="flex items-center px-1 gap-1") {
                                        img(src=cloned.logo_url, class="w-6 h-6") {}
                                        span(class="text-sm") { (shortened_address) }
                                    }
                                    div(class="flex flex-wrap items-center gap-4") {
                                        div(class=button_row_class, on:click=move |_| handle_update(Update::NotifyFunding)) {
                                            span(class=(if wallet.get().notify_funding { selected_class } else { unselected_class })) {}
                                            span(class="text-sm") { "Funding" }
                                        }
                                        div(class=button_row_class, on:click=move |_| handle_update(Update::NotifyStaking)) {
                                            span(class=(if wallet.get().notify_staking { selected_class } else { unselected_class })) {}
                                            span(class="text-sm") { "Staking" }
                                        }
                                        div(class=button_row_class, on:click=move |_| handle_update(Update::NotifyGovVotingReminder)) {
                                            span(class=(if wallet.get().notify_gov_voting_reminder { selected_class } else { unselected_class })) {}
                                            span(class="text-sm truncate") { "Governance reminders" }
                                        }
                                    }
                                }
                            }
                        },
                    )
                }
            }
        }
    }
}
