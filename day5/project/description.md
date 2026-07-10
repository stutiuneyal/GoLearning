# Project: E-Commerce Order Management System

Build a command-line **E-Commerce Order Management System** that simulates the process of managing customers, products, shopping carts, orders, and payments for
an online shopping platform.

The application should maintain customer information, including separate billing and shipping addresses. A customer can purchase one or more products,
where each product contains details such as product ID, name, price, and available stock. An order should consist of a customer and a collection of cart items,
with each cart item representing a product and its purchased quantity. Design these relationships using composition.

The application must support the following operations:

1. Create billing and shipping addresses.
2. Create customers using a constructor function.
3. Update a customer's email address after validating the input.
4. Display complete customer details, including billing and shipping addresses.
5. Create multiple products and maintain them in a product catalog.
6. Display all available products.
7. Add one or more products to an order with the required quantities.
8. Calculate the total value of an order based on the products and quantities purchased.
9. Define order statuses (**Created**, **Paid**, **Shipped**, **Delivered**, and **Cancelled**) using `iota`.
10. Define payment statuses (**Pending**, **Successful**, and **Failed**) using `iota`.
11. Support multiple payment methods, such as **Card Payment** and **UPI Payment**, through a common interface.
12. Process the payment for an order and update the order status based on the payment result.
13. Generate and display a complete invoice containing customer information, shipping and billing addresses, purchased products, quantities, payment details,
order status, and total amount.
14. Implement meaningful string representations for appropriate objects to improve console output.

Organize the application using constructors, methods, composition, interfaces, polymorphism, encapsulation, enumerations, slices, maps,
variadic functions where appropriate, and error handling. The application should produce a clean and well-formatted command-line output
that simulates a simple e-commerce order processing system.
