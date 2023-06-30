use std::fmt::Debug;

use crate::config::keys;
use crate::types::protobuf::grpc::user_setup_service_client::UserSetupServiceClient;
use grpc_web_client::Client;
use tonic::metadata::MetadataValue;
use tonic::Request;

use crate::services::auth::AuthService;
use crate::types::protobuf::grpc::auth_service_client::AuthServiceClient;
use crate::types::protobuf::grpc::event_service_client::EventServiceClient;
use crate::types::protobuf::grpc::user_service_client::UserServiceClient;

#[derive(Debug, Clone)]
pub struct GrpcClient {
    endpoint_url: String,
    auth_manager: AuthService,
}

impl Default for GrpcClient {
    fn default() -> Self {
        Self::new()
    }
}

impl GrpcClient {
    pub fn new() -> Self {
        Self {
            endpoint_url: keys::GRPC_WEB_ENDPOINT_URL.to_string(),
            auth_manager: AuthService::new(),
        }
    }

    pub fn create_request<T>(&self, message: T) -> Request<T> {
        let token = self.auth_manager.get_access_token();
        let mut req = Request::new(message);
        if let Ok(token) = token {
            req.metadata_mut().insert(
                "authorization",
                MetadataValue::try_from(token).unwrap_or_else(|_| MetadataValue::from_static("")),
            );
        }
        req
    }

    pub fn get_auth_service(&self) -> AuthServiceClient<Client> {
        AuthServiceClient::new(Client::new(self.endpoint_url.clone()))
    }

    pub fn get_user_service(&self) -> UserServiceClient<Client> {
        UserServiceClient::new(Client::new(self.endpoint_url.clone()))
    }

    pub fn get_user_setup_service(&self) -> UserSetupServiceClient<Client> {
        UserSetupServiceClient::new(Client::new(self.endpoint_url.clone()))
    }

    pub fn get_event_service(&self) -> EventServiceClient<Client> {
        EventServiceClient::new(Client::new(self.endpoint_url.clone()))
    }
}
