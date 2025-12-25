package models

type ShippingMethodModel struct {
	Name          string  `json:"name" validate:"required,min=3,max=255"`
	Description   *string `json:"description" validate:"omitempty,max=1000"`
	Price         float64 `json:"price" validate:"required,gte=0"`
	EstimatedDays *int    `json:"estimated_days" validate:"omitempty,gte=0"`
	Active        bool    `json:"active" validate:"required"`
}

type ShippingZoneModel struct {
	Name        string   `json:"name" validate:"required,min=3,max=255"`
	Description *string  `json:"description" validate:"omitempty,max=1000"`
	Countries   []string `json:"countries" validate:"required,dive,iso3166_1_alpha2"`
}

type ShippingRateModel struct {
	ShippingMethodID string  `json:"shipping_method_id" validate:"required,uuid4"`
	ShippingZoneID   string  `json:"shipping_zone_id" validate:"required,uuid4"`
	Price            float64 `json:"price" validate:"required,gte=0"`
}

type UpdateShippingRateModel struct {
	Price *float64 `json:"price" validate:"omitnil,omitempty,gte=0"`
}

type ShipmentModel struct {
	Name         string  `json:"name" validate:"required,min=3,max=255"`
	Phone        string  `json:"phone" validate:"required,min=10,max=15"`
	Address      string  `json:"address" validate:"required,min=10,max=500"`
	City         string  `json:"city" validate:"required,min=2,max=100"`
	State        string  `json:"state" validate:"required,min=2,max=100"`
	Country      string  `json:"country" validate:"required,min=2,max=100"`
	PostalCode   string  `json:"postal_code" validate:"required,min=4,max=20"`
	Instructions *string `json:"instructions" validate:"omitempty,max=1000"`
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
