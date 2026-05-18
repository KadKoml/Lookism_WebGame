package grpc

import (
	"context"

	pb "github.com/kozie/lookism-rpg/api/proto/notification"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	// In a real app we would have a usecase here: uc *usecase.NotificationUsecase
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

// History

func (h *NotificationHandler) GetNotificationHistory(ctx context.Context, req *pb.UserRequest) (*pb.HistoryResponse, error) {
	// Stub
	return &pb.HistoryResponse{}, nil
}

func (h *NotificationHandler) ClearNotificationHistory(ctx context.Context, req *pb.UserRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{Success: true}, nil
}

func (h *NotificationHandler) DeleteNotification(ctx context.Context, req *pb.DeleteRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{Success: true}, nil
}

func (h *NotificationHandler) MarkAsRead(ctx context.Context, req *pb.MarkReadRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{Success: true}, nil
}

// Settings

func (h *NotificationHandler) GetSettings(ctx context.Context, req *pb.UserRequest) (*pb.SettingsResponse, error) {
	return &pb.SettingsResponse{EmailEnabled: true}, nil
}

func (h *NotificationHandler) UpdateSettings(ctx context.Context, req *pb.UpdateSettingsRequest) (*pb.SettingsResponse, error) {
	return &pb.SettingsResponse{
		EmailEnabled: req.EmailEnabled,
		PushEnabled:  req.PushEnabled,
		SmsEnabled:   req.SmsEnabled,
	}, nil
}

func (h *NotificationHandler) OptOutEmails(ctx context.Context, req *pb.UserRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{Success: true}, nil
}

func (h *NotificationHandler) OptInEmails(ctx context.Context, req *pb.UserRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{Success: true}, nil
}

// Templates

func (h *NotificationHandler) GetEmailTemplates(ctx context.Context, req *pb.EmptyRequest) (*pb.TemplatesResponse, error) {
	return &pb.TemplatesResponse{}, nil
}

func (h *NotificationHandler) PreviewTemplate(ctx context.Context, req *pb.PreviewRequest) (*pb.PreviewResponse, error) {
	return &pb.PreviewResponse{PreviewHtml: "<p>Preview</p>"}, nil
}

// Manual dispatch (admin)

func (h *NotificationHandler) SendCustomNotification(ctx context.Context, req *pb.SendCustomRequest) (*pb.EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCustomNotification not implemented")
}

func (h *NotificationHandler) BroadcastNotification(ctx context.Context, req *pb.BroadcastRequest) (*pb.EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastNotification not implemented")
}
