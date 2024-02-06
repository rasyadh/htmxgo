# HTMX-GO

## Overview

HTMX-GO is a web application project that demonstrates the integration of Go (using Echo framework), HTMX, and Tailwind CSS. The project showcases the seamless integration of htmx capabilities for dynamic client-server communication while leveraging the power of Go templates for server-side rendering.

## Technologies Used

- [Go Echo](https://github.com/labstack/echo): A fast and minimalist Go web framework.
- [HTMX](https://htmx.org/): A lightweight JavaScript library for making AJAX-driven websites easy to create.
- [Tailwind CSS](https://tailwindcss.com/): A utility-first CSS framework for rapidly building custom designs.
- [DaisyUI](https://daisyui.com/): A UI library for Tailwind CSS, extending its capabilities with additional components and utilities.
- [Air](https://github.com/cosmtrek/air): A live reload tool for Go applications.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/): Ensure that Go is installed on your machine.
- [Node.js](https://nodejs.org/): Required for managing the frontend dependencies.
- [Tailwind CSS CLI](https://tailwindcss.com/docs/installation#using-tailwind-cli): Install Tailwind CSS CLI for building styles.
- [Air](https://github.com/cosmtrek/air): Install Air for live reloading Go applications. You can install it using:

  ```bash
  go get -u github.com/cosmtrek/air
  ```

### Installation

1. Clone the repository:

   ```bash
   git clone https://gitlab.com/your-username/htmx-go.git
   cd htmx-go
   ```

2. Install Go dependencies:

    ```bash
    go mod tidy
    ```

3. Run the application with Air:

    ```bash
    air
    ```

    Air will handle live reloading as you make changes to your Go code.
