# Food Delivery workflow

1. User creates order via API Gateway

   - Order Service validates request and creates order (status: PENDING)

2. Order Service publishes OrderCreated event to message broker

3. Payment Service consumes OrderCreated event

   - Processes payment
   - Publishes PaymentCompleted or PaymentFailed event

4. Order Service consumes PaymentCompleted event

   - Updates order status to CONFIRMED
   - Publishes OrderConfirmed event

5. Food/Restaurant Service consumes OrderConfirmed event

   - Checks food availability
   - Notifies restaurant/kitchen
   - Publishes OrderAccepted event (or OrderRejected if unavailable)

6. Order Service consumes OrderAccepted event

   - Updates order status to PREPARING

7. Delivery Service consumes OrderAccepted event

   - Assigns delivery driver
   - Publishes DriverAssigned event

8. Food Service publishes FoodReady event when preparation complete

9. Order Service updates status to READY_FOR_PICKUP

10. Delivery Service picks up order

    - Publishes OrderPickedUp event
    - Order status → OUT_FOR_DELIVERY

11. Delivery Service completes delivery

    - Publishes OrderDelivered event
    - Order status → DELIVERED

12. Notification Service (optional) consumes events and notifies user at each step
