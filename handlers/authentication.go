package handlers

// import (
//     "time"

//     "golang.org/x/crypto/bcrypt"
//     "github.com/dgrijalva/jwt-go"
//     "github.com/gofiber/fiber/v2"
//     "github.com/nikhil478/function-wale-server/models"
//     "gorm.io/gorm"
// )

// var jwtKey = []byte("your_secret_key")

// type Credentials struct {
//     Email    string `json:"email"`
//     Password string `json:"password"`
// }

// type Claims struct {
//     Email string `json:"email"`
//     Role  string `json:"role"`
//     jwt.StandardClaims
// }

// // SignUp handles the user registration
// // @Summary Sign up a new user
// // @Tags authentication
// // @Accept json
// // @Produce json
// // @Param credentials body Credentials true "User Credentials"
// // @Success 201 {object} map[string]interface{}
// // @Failure 400 {object} map[string]string
// // @Failure 409 {object} map[string]string
// // @Router /api/signup [post]
// func SignUp(db *gorm.DB) fiber.Handler {
//     return func(ctx *fiber.Ctx) error {
//         var creds Credentials
//         if err := ctx.BodyParser(&creds); err != nil {
//             ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid request"})
//             return err
//         }

//         if creds.Email == "" || creds.Password == "" {
//             ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Email and Password required"})
//             return nil
//         }

//         var existingUser models.User
//         if err := db.Where("email = ?", creds.Email).First(&existingUser).Error; err == nil {
//             ctx.Status(fiber.StatusConflict).JSON(&fiber.Map{"message": "User already exists"})
//             return nil
//         }

//         hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
//         if err != nil {
//             ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not hash password"})
//             return err
//         }

//         user := models.User{
//             Email:    creds.Email,
//             Password: string(hashedPassword),
//             Role:     "user",
//         }

//         if err := db.Create(&user).Error; err != nil {
//             ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not create user"})
//             return err
//         }

//         expirationTime := time.Now().Add(15 * time.Minute)
//         claims := &Claims{
//             Email: user.Email,
//             Role:  user.Role,
//             StandardClaims: jwt.StandardClaims{
//                 ExpiresAt: expirationTime.Unix(),
//             },
//         }

//         token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//         tokenString, err := token.SignedString(jwtKey)
//         if err != nil {
//             ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not create token"})
//             return err
//         }

//         ctx.Status(fiber.StatusCreated).JSON(&fiber.Map{
//             "message": "User created successfully",
//             "token":   tokenString,
//         })
//         return nil
//     }
// }

// // Login handles the user login and JWT token issuance
// // @Summary Log in a user
// // @Tags authentication
// // @Accept json
// // @Produce json
// // @Param credentials body Credentials true "User Credentials"
// // @Success 200 {object} map[string]interface{}
// // @Failure 401 {object} map[string]string
// // @Router /api/login [post]
// func Login(db *gorm.DB) fiber.Handler {
//     return func(ctx *fiber.Ctx) error {
//         var creds Credentials
//         if err := ctx.BodyParser(&creds); err != nil {
//             ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid request"})
//             return err
//         }

//         var user models.User
//         if err := db.Where("email = ?", creds.Email).First(&user).Error; err != nil {
//             ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid email or password"})
//             return err
//         }

//         if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
//             ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid email or password"})
//             return nil
//         }

//         expirationTime := time.Now().Add(15 * time.Minute)
//         claims := &Claims{
//             Email: user.Email,
//             Role:  user.Role,
//             StandardClaims: jwt.StandardClaims{
//                 ExpiresAt: expirationTime.Unix(),
//             },
//         }

//         token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//         tokenString, err := token.SignedString(jwtKey)
//         if err != nil {
//             ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not create token"})
//             return err
//         }

//         ctx.JSON(&fiber.Map{
//             "message": "Login successful",
//             "token":   tokenString,
//         })
//         return nil
//     }
// }

// // JWTMiddleware is the middleware for protecting routes
// func JWTMiddleware() fiber.Handler {
//     return func(ctx *fiber.Ctx) error {
//         tokenString := ctx.Get("Authorization")
//         if tokenString == "" {
//             ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Missing or invalid token"})
//             return nil
//         }

//         claims := &Claims{}
//         token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
//             return jwtKey, nil
//         })
//         if err != nil || !token.Valid {
//             ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid token"})
//             return nil
//         }

//         ctx.Locals("user", claims)
//         return ctx.Next()
//     }
// }
