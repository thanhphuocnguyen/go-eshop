package models

type ShippingMethodModel struct {
	Name            string  `json:"name" validate:"required,min=3,max=255"`
	Description     *string `json:"description" validate:"omitempty,max=1000"`
	Price           float64 `json:"price" validate:"required,gte=0"`
	EstimatedDays   *int    `json:"estimated_days" validate:"omitempty,gte=0"`
	RequiresAddress bool    `json:"requires_address" validate:"required"`
	Active          bool    `json:"active" validate:"required"`
}

type ShippingZoneModel struct {
	Name        string   `json:"name" validate:"required,min=3,max=255"`
	Description *string  `json:"description" validate:"omitempty,max=1000"`
	Countries   []string `json:"countries" validate:"required,dive,iso3166_1_alpha2"`
}

type ShippingRateModel struct {
	ShippingMethodID      string   `json:"shipping_method_id" validate:"required,uuid4"`
	ShippingZoneID        string   `json:"shipping_zone_id" validate:"required,uuid4"`
	Price                 float64  `json:"price" validate:"required,gte=0"`
	MinOrderAmount        *float64 `json:"min_order_amount" validate:"omitempty,gte=0"`
	MaxOrderAmount        *float64 `json:"max_order_amount" validate:"omitempty,gte=0"`
	FreeShippingThreshold *float64 `json:"free_shipping_threshold" validate:"omitempty,gte=0"`
	IsActive              bool     `json:"is_active" validate:"required"`
	Name                  string   `json:"name" validate:"omitempty,min=3,max=255"`
}

type UpdateShippingRateModel struct {
	Price *float64 `json:"price" validate:"omitnil,omitempty,gte=0"`
}

type ShipmentModel struct {
	OrderID          string  `json:"order_id" validate:"required,uuid4"`
	Status           string  `json:"status" validate:"required,oneof='pending' 'shipped' 'delivered' 'cancelled'"`
	ShippedAt        *string `json:"shipped_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	DeliveredAt      *string `json:"delivered_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	TrackingNumber   *string `json:"tracking_number" validate:"omitempty,min=3,max=100"`
	TrackingUrl      *string `json:"tracking_url" validate:"omitempty,url"`
	ShippingProvider string  `json:"shipping_provider" validate:"omitempty,min=2,max=100"`
	ShippingNotes    string  `json:"shipping_notes" validate:"omitempty,max=1000"`
}

type UpdateShipmentModel struct {
	Name         *string `json:"name" validate:"omitnil,omitempty,min=3,max=255"`
	Phone        *string `json:"phone" validate:"omitnil,omitempty,min=10,max=15"`
	Address      *string `json:"address" validate:"omitnil,omitempty,min=10,max=500"`
	City         *string `json:"city" validate:"omitnil,omitempty,min=2,max=100"`
	State        *string `json:"state" validate:"omitnil,omitempty,min=2,max=100"`
	Country      *string `json:"country" validate:"omitnil,omitempty,min=2,max=100"`
	PostalCode   *string `json:"postal_code" validate:"omitnil,omitempty,min=4,max=20"`
	Instructions *string `json:"instructions" validate:"omitempty,max=1000"`
}

type ShipmentStatusModel struct {
	Status string `json:"status" validate:"required,oneof='pending' 'shipped' 'delivered' 'cancelled'"`
}

type ShipmentTrackingModel struct {
	TrackingNumber *string `json:"tracking_number" validate:"omitnil,omitempty,min=3,max=100"`
	Carrier        *string `json:"carrier" validate:"omitnil,omitempty,min=2,max=100"`
}
