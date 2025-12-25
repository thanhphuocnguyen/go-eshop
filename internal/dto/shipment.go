package dto

type ShipmentListItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Price       float64 `json:"price"`
	EstimatedAt *string `json:"estimatedAt"`
	IsActive    bool    `json:"isActive"`
}

type ShipmentDetail struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    *string  `json:"description"`
	Price          float64  `json:"price"`
	EstimatedAt    *string  `json:"estimatedAt"`
	IsActive       bool     `json:"isActive"`
	Countries      []string `json:"countries"`
	TrackingNumber *string  `json:"trackingNumber,omitempty"`
	Carrier        *string  `json:"carrier,omitempty"`
}

type ShipmentAddress struct {
	Name         string  `json:"name"`
	Phone        string  `json:"phone"`
	Address      string  `json:"address"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	Country      string  `json:"country"`
	PostalCode   string  `json:"postalCode"`
	Instructions *string `json:"instructions,omitempty"`
}

type ShipmentTrackingInfo struct {
	TrackingNumber string `json:"trackingNumber"`
	Carrier        string `json:"carrier"`
	Status         string `json:"status"`
}

type ShipmentUpdateTracking struct {
	TrackingNumber *string `json:"trackingNumber,omitempty"`
	Carrier        *string `json:"carrier,omitempty"`
}

type ShipmentUpdateStatus struct {
	Status string `json:"status"`
}
type ShipmentAssignTracking struct {
	TrackingNumber string `json:"trackingNumber"`
	Carrier        string `json:"carrier"`
}
type ShipmentAssignStatus struct {
	Status string `json:"status"`
}
type ShipmentCreateModel struct {
	Name         string  `json:"name" validate:"required,min=3,max=255"`
	Phone        string  `json:"phone" validate:"required,min=10,max=15"`
	Address      string  `json:"address" validate:"required,min=10,max=500"`
	City         string  `json:"city" validate:"required,min=2,max=100"`
	State        string  `json:"state" validate:"required,min=2,max=100"`
	Country      string  `json:"country" validate:"required,min=2,max=100"`
	PostalCode   string  `json:"postal_code" validate:"required,min=4,max=20"`
	Instructions *string `json:"instructions" validate:"omitempty,max=1000"`
}

type ShipmentUpdateModel struct {
	Name         *string `json:"name" validate:"omitnil,omitempty,min=3,max=255"`
	Phone        *string `json:"phone" validate:"omitnil,omitempty,min=10,max=15"`
	Address      *string `json:"address" validate:"omitnil,omitempty,min=10,max=500"`
	City         *string `json:"city" validate:"omitnil,omitempty,min=2,max=100"`
	State        *string `json:"state" validate:"omitnil,omitempty,min=2,max=100"`
	Country      *string `json:"country" validate:"omitnil,omitempty,min=2,max=100"`
	PostalCode   *string `json:"postal_code" validate:"omitnil,omitempty,min=4,max=20"`
	Instructions *string `json:"instructions" validate:"omitempty,max=1000"`
}
type ShipmentStatusUpdateModel struct {
	Status string `json:"status" validate:"required,oneof='pending' 'shipped' 'delivered' 'cancelled'"`
}

type ShipmentTrackingUpdateModel struct {
	TrackingNumber *string `json:"tracking_number" validate:"omitnil,omitempty,min=3,max=100"`
	Carrier        *string `json:"carrier" validate:"omitnil,omitempty,min=2,max=100"`
}
