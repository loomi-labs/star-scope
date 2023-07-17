pub mod grpc_auth {
    tonic::include_proto!("starscope.grpc.auth");
}

pub mod grpc_event {
    tonic::include_proto!("starscope.grpc.event");
}

pub mod grpc_indexer {
    tonic::include_proto!("starscope.grpc.indexer");
}

pub mod grpc_settings {
    tonic::include_proto!("starscope.grpc.settings");
}

pub mod grpc_user {
    tonic::include_proto!("starscope.grpc.user");
}