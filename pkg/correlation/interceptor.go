/*
Copyright 2021 Portworx

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package correlation

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ContextInterceptor represents a correlation interceptor
type ContextInterceptor struct {
	Origin Component
}

// ContextUnaryInterceptor creates a gRPC interceptor for adding
// correlation ID to each request
func (ci *ContextInterceptor) ContextUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// Get request context from gRPC metadata
		rc := RequestContextFromContextMetadata(md)

		// Only add to context if an ID exists
		if len(rc.ID) > 0 {
			ctx = context.WithValue(ctx, ContextKey, rc)
		}
	}
	ctx = WithCorrelationContext(ctx, ci.Origin)

	return handler(ctx, req)
}
