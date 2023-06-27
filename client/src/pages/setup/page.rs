use sycamore::{prelude::*, view};
use sycamore::futures::spawn_local_scoped;
use tonic::Status;
use crate::components::button::{SolidButton, OutlineButton, ColorScheme};

use crate::components::loading::LoadingSpinner;
use crate::components::messages::{create_error_msg_from_status, create_message};
use crate::Services;
use crate::components::navigation::Header;
use crate::types::protobuf::grpc::{finish_step_request, FinishStepRequest, StepOneRequest, StepOneResponse, StepResponse, StepTwoResponse, StepThreeResponse, GetStepRequest, get_step_request, StepThreeRequest, StepTwoRequest, StepFourRequest, StepFourResponse, StepFiveRequest, StepFiveResponse};
use crate::types::protobuf::grpc::step_response::Step;
use crate::types::protobuf::grpc::step_response::Step::{StepFive, StepFour, StepOne, StepThree, StepTwo};

#[derive(Debug, Clone)]
pub struct SetupState {
    pub step: RcSignal<Option<Step>>,
}

impl SetupState {
    pub fn new() -> Self {
        Self {
            step: create_rc_signal(None),
        }
    }
}

const TITLE_CLASS: &str = "text-4xl font-bold my-4";
const SUBTITLE_CLASS: &str = "text-2xl font-semibold my-2";
const DESCRIPTION_CLASS: &str = "dark:text-purple-600";
const BUTTON_ROW_CLASS: &str = "flex justify-center space-x-4";


#[component(inline_props)]
pub fn StepOneComponent<G: Html>(cx: Scope, step: StepOneResponse) -> View<G> {
    let handle_click = move |is_validator| {
        spawn_local_scoped(cx, async move {
            let finish_step = FinishStepRequest {
                go_to_next_step: true,
                step: Some(finish_step_request::Step::StepOne(StepOneRequest {
                    is_validator,
                })),
            };
            update_step(cx, finish_step).await;
        });
    };

    view! {cx,
        h1(class=TITLE_CLASS) {"Welcome to Star Scope!"}
        p(class=DESCRIPTION_CLASS) {"We deliver quick and effortless notifications about your Cosmos ecosystem activities."}
        h2(class=SUBTITLE_CLASS) {"Are you a validator?"}
        div(class=BUTTON_ROW_CLASS) {
            SolidButton(on_click=move || handle_click(true), color=ColorScheme::Subtle) {"Yes"}
            SolidButton(on_click=move || handle_click(false), color=ColorScheme::Subtle) {"No"}
        }
    }
}

#[component(inline_props)]
pub fn StepTwoComponent<G: Html>(cx: Scope, step: StepTwoResponse) -> View<G> {
    let handle_click = move |go_to_next_step| {
        spawn_local_scoped(cx, async move {
            let finish_step = FinishStepRequest {
                go_to_next_step,
                step: Some(finish_step_request::Step::StepTwo(StepTwoRequest {
                    validator_ids: vec![],
                })),
            };
            update_step(cx, finish_step).await;
        });
    };

    view! {cx,
        h2(class=TITLE_CLASS) {"Choose your validator (s)"}
        p(class=DESCRIPTION_CLASS) {"You will receive reminders to vote on governance proposals from the validators you've selected."}
        div(class=BUTTON_ROW_CLASS) {
            OutlineButton(on_click=move || handle_click(false)) {"Back"}
            SolidButton(on_click=move || handle_click(true)) {"Next"}
        }
    }
}

#[component(inline_props)]
pub fn StepThreeComponent<G: Html>(cx: Scope, step: StepThreeResponse) -> View<G> {
    let handle_click = move |go_to_next_step| {
        spawn_local_scoped(cx, async move {
            let finish_step = FinishStepRequest {
                go_to_next_step,
                step: Some(finish_step_request::Step::StepThree(StepThreeRequest {wallet_addresses: vec![]})),
            };
            update_step(cx, finish_step).await;
        });
    };

    view! {cx,
        p(class=DESCRIPTION_CLASS) {"Add your wallet(s)"}
        div(class=BUTTON_ROW_CLASS) {
            OutlineButton(on_click=move || handle_click(false)) {"Back"}
            SolidButton(on_click=move || handle_click(true)) {"Next"}
        }
    }
}

#[component(inline_props)]
pub fn StepFourComponent<G: Html>(cx: Scope, step: StepFourResponse) -> View<G> {
    let handle_click = move |go_to_next_step| {
        spawn_local_scoped(cx, async move {
            let finish_step = FinishStepRequest {
                go_to_next_step,
                step: Some(finish_step_request::Step::StepFour(StepFourRequest {
                    notify_funding: false,
                    ..Default::default()
                })),
            };
            update_step(cx, finish_step).await;
        });
    };

    view! {cx,
        p(class=DESCRIPTION_CLASS) {"Choose your Notifications"}
        div(class=BUTTON_ROW_CLASS) {
            OutlineButton(on_click=move || handle_click(false)) {"Back"}
            SolidButton(on_click=move || handle_click(true)) {"Next"}
        }
    }
}

#[component(inline_props)]
pub fn StepFiveComponent<G: Html>(cx: Scope, step: StepFiveResponse) -> View<G> {
    let handle_click = move |go_to_next_step| {
        spawn_local_scoped(cx, async move {
            let finish_step = FinishStepRequest {
                go_to_next_step,
                step: Some(finish_step_request::Step::StepFive(StepFiveRequest {
                    ..Default::default()
                })),
            };
            update_step(cx, finish_step).await;
        });
    };

    view! {cx,
        p(class=DESCRIPTION_CLASS) {"Choose where to receive notifications"}
        div(class=BUTTON_ROW_CLASS) {
            OutlineButton(on_click=move || handle_click(false)) {"Back"}
            SolidButton(on_click=move || handle_click(true)) {"Finish"}
        }
    }
}

async fn update_step(cx: Scope<'_>, finish_step: FinishStepRequest) {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(finish_step);
    let response = services
        .grpc_client
        .get_user_setup_service()
        .finish_step(request)
        .await
        .map(|res| res.into_inner());
    handle_setup_step_response(cx, response)
}

fn handle_setup_step_response(cx: Scope, response: Result<StepResponse, Status>) {
    if let Ok(response) = response {
        let setup_state = use_context::<SetupState>(cx);
        match response.step {
            Some(step) => {
                setup_state.step.set(Some(step));
            }
            None => create_message(cx, "Setup step not found", "Setup step was not found", crate::InfoLevel::Error),
        }
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

async fn query_setup_step(cx: Scope<'_>, requestedStep: Option<get_step_request::Step>) {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(GetStepRequest {step: requestedStep.unwrap_or_else(|| get_step_request::Step::CurrentStep) as i32});
    let response = services
        .grpc_client
        .get_user_setup_service()
        .get_step(request)
        .await
        .map(|res| res.into_inner());
    handle_setup_step_response(cx, response)
}


#[component]
pub fn Setup<G: Html>(cx: Scope) -> View<G> {
    let setup_state = SetupState::new();
    provide_context(cx, setup_state.clone());

    spawn_local_scoped(cx, async move {
        query_setup_step(cx, None).await;
    });

    view! {cx,
        (if let Some(step) = setup_state.step.get().as_ref() {
            let child = match step {
                StepOne(s) => view! {cx, StepOneComponent(step=s.clone())},
                StepTwo(s) => view! {cx, StepTwoComponent(step=s.clone())},
                StepThree(s) => view! {cx, StepThreeComponent(step=s.clone())},
                StepFour(s) => view! {cx, StepFourComponent(step=s.clone())},
                StepFive(s) => view! {cx, StepFiveComponent(step=s.clone())},
            };
            view! {
                cx,
                div(class="min-h-[100svh] flex justify-center items-center flex-auto flex-shrink-0") {
                    div(class="flex flex-col lg:max-w-screen-lg xl:max-w-screen-xl h-full w-full") {
                        Header{}
                        div(class="flex flex-col w-full min-h-[100svh] text-center p-4") {
                            div(class="flex flex-col flex-auto justify-center items-center") {
                                (child)
                            }
                        }
                    }
                }
            }
        } else {
            view! {cx,
                LoadingSpinner {}
            }
        })
    }
}
