package v1

import (
	"fmt"
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/admin"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/auth"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/cart"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/catalog"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/checkout"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/content"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/importer"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/lead"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/media"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/notification"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/orders"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/profile"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/reviews"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/user"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/wishlist"
	cartmodule "github.com/evrone/go-clean-template/internal/module/cart"
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	"github.com/evrone/go-clean-template/internal/module/catalog/catalogbase"
	importermodule "github.com/evrone/go-clean-template/internal/module/importer"
	leadmodule "github.com/evrone/go-clean-template/internal/module/lead"
	mediamodule "github.com/evrone/go-clean-template/internal/module/media"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
)

var _ openapi.StrictServerInterface = (*Server)(nil)

// Server is the delivery layer composition root for OpenAPI generated handlers.
type Server struct {
	cfg         *config.Config
	userUC      usermodule.UserUseCase
	catalogUC   catalogmodule.CatalogUseCase
	catalogBase catalogbase.UseCase
	authUC      usermodule.AuthUseCase
	checkoutUC  ordermodule.CheckoutUseCase
	ordersUC    ordermodule.OrdersUseCase
	profileUC   usermodule.ProfileUseCase
	adminUC     usermodule.AdminUseCase
	wishlistUC  wishlistmodule.WishlistUseCase
	notifyUC    notificationmodule.NotificationCenterUseCase
	mediaUC     mediamodule.MediaUseCase
	homepageUC  contentmodule.HomepageUseCase
	cartUC      cartmodule.CartUseCase
	leadUC      leadmodule.LeadUseCase
	contentUC   contentmodule.ContentUseCase
	importerUC  importermodule.ImporterUseCase
	jwtManager  *jwt.Manager
	logger      logger.Interface
}

// NewServer constructs the OpenAPI Server composition root.
func NewServer(
	cfg *config.Config,
	userUC usermodule.UserUseCase,
	catalogUC catalogmodule.CatalogUseCase,
	catalogBase catalogbase.UseCase,
	authUC usermodule.AuthUseCase,
	checkoutUC ordermodule.CheckoutUseCase,
	ordersUC ordermodule.OrdersUseCase,
	profileUC usermodule.ProfileUseCase,
	adminUC usermodule.AdminUseCase,
	wishlistUC wishlistmodule.WishlistUseCase,
	notifyUC notificationmodule.NotificationCenterUseCase,
	mediaUC mediamodule.MediaUseCase,
	homepageUC contentmodule.HomepageUseCase,
	cartUC cartmodule.CartUseCase,
	leadUC leadmodule.LeadUseCase,
	contentUC contentmodule.ContentUseCase,
	importerUC importermodule.ImporterUseCase,
	jwtManager *jwt.Manager,
	l logger.Interface,
) *Server {
	return &Server{
		cfg:         cfg,
		userUC:      userUC,
		catalogUC:   catalogUC,
		catalogBase: catalogBase,
		authUC:      authUC,
		checkoutUC:  checkoutUC,
		ordersUC:    ordersUC,
		profileUC:   profileUC,
		adminUC:     adminUC,
		wishlistUC:  wishlistUC,
		notifyUC:    notifyUC,
		mediaUC:     mediaUC,
		homepageUC:  homepageUC,
		cartUC:      cartUC,
		leadUC:      leadUC,
		contentUC:   contentUC,
		importerUC:  importerUC,
		jwtManager:  jwtManager,
		logger:      l,
	}
}

// Handler instantiations for vertical capabilities

func (s *Server) adminHandler() *admin.Handler {
	return admin.NewHandler(s.adminUC, s.logger)
}

func (s *Server) authHandler() *auth.Handler {
	return auth.NewHandler(s.authUC, s.userUC, s.logger)
}

func (s *Server) cartHandler() *cart.Handler {
	return cart.NewHandler(s.cartUC, s.logger)
}

func (s *Server) catalogHandler() *catalog.Handler {
	return catalog.NewHandler(s.catalogUC, s.catalogBase, s.logger)
}

func (s *Server) checkoutHandler() *checkout.Handler {
	return checkout.NewHandler(s.checkoutUC, s.logger)
}

func (s *Server) contentHandler() *content.Handler {
	return content.NewHandler(s.contentUC, s.homepageUC, s.logger)
}

func (s *Server) importerHandler() *importer.Handler {
	return importer.NewHandler(s.importerUC, s.logger)
}

func (s *Server) leadHandler() *lead.Handler {
	return lead.NewHandler(s.leadUC, s.logger)
}

func (s *Server) mediaHandler() *media.Handler {
	return media.NewHandler(s.mediaUC, s.logger)
}

func (s *Server) notificationHandler() *notification.Handler {
	return notification.NewHandler(s.notifyUC, s.logger)
}

func (s *Server) ordersHandler() *orders.Handler {
	return orders.NewHandler(s.ordersUC, s.logger)
}

func (s *Server) profileHandler() *profile.Handler {
	return profile.NewHandler(s.profileUC, s.logger)
}

func (s *Server) userHandler() *user.Handler {
	return user.NewHandler(s.userUC, s.logger)
}

func (s *Server) wishlistHandler() *wishlist.Handler {
	return wishlist.NewHandler(s.wishlistUC, s.logger)
}

func (s *Server) reviewsHandler() *reviews.Handler {
	return reviews.NewHandler(s.logger)
}

// -----------------------------------------------------------------------------
// OpenAPI StrictServerInterface Delegation Methods
// -----------------------------------------------------------------------------

// Admin Capability
func (s *Server) GetAdminDashboardStats(ctx context.Context, request openapi.GetAdminDashboardStatsRequestObject) (openapi.GetAdminDashboardStatsResponseObject, error) {
	return s.adminHandler().GetAdminDashboardStats(ctx, request)
}

func (s *Server) ListAdminAuditLogs(ctx context.Context, request openapi.ListAdminAuditLogsRequestObject) (openapi.ListAdminAuditLogsResponseObject, error) {
	return s.adminHandler().ListAdminAuditLogs(ctx, request)
}

// Admin Orders
func (s *Server) AdminListOrders(ctx context.Context, request openapi.AdminListOrdersRequestObject) (openapi.AdminListOrdersResponseObject, error) {
	return s.adminHandler().AdminListOrders(ctx, request)
}

func (s *Server) AdminGetOrder(ctx context.Context, request openapi.AdminGetOrderRequestObject) (openapi.AdminGetOrderResponseObject, error) {
	return s.adminHandler().AdminGetOrder(ctx, request)
}

func (s *Server) AdminUpdateOrder(ctx context.Context, request openapi.AdminUpdateOrderRequestObject) (openapi.AdminUpdateOrderResponseObject, error) {
	return s.adminHandler().AdminUpdateOrder(ctx, request)
}

func (s *Server) AdminGetCollect(ctx context.Context, request openapi.AdminGetCollectRequestObject) (openapi.AdminGetCollectResponseObject, error) {
	return s.adminHandler().AdminGetCollect(ctx, request)
}

// Admin Refunds
func (s *Server) AdminListRefunds(ctx context.Context, request openapi.AdminListRefundsRequestObject) (openapi.AdminListRefundsResponseObject, error) {
	return s.adminHandler().AdminListRefunds(ctx, request)
}

func (s *Server) AdminApproveRefund(ctx context.Context, request openapi.AdminApproveRefundRequestObject) (openapi.AdminApproveRefundResponseObject, error) {
	return s.adminHandler().AdminApproveRefund(ctx, request)
}

func (s *Server) AdminRejectRefund(ctx context.Context, request openapi.AdminRejectRefundRequestObject) (openapi.AdminRejectRefundResponseObject, error) {
	return s.adminHandler().AdminRejectRefund(ctx, request)
}

// Admin Reviews
func (s *Server) AdminListReviews(ctx context.Context, request openapi.AdminListReviewsRequestObject) (openapi.AdminListReviewsResponseObject, error) {
	return s.adminHandler().AdminListReviews(ctx, request)
}

func (s *Server) AdminPublishSelectedReviews(ctx context.Context, request openapi.AdminPublishSelectedReviewsRequestObject) (openapi.AdminPublishSelectedReviewsResponseObject, error) {
	return s.adminHandler().AdminPublishSelectedReviews(ctx, request)
}

func (s *Server) AdminDeleteReview(ctx context.Context, request openapi.AdminDeleteReviewRequestObject) (openapi.AdminDeleteReviewResponseObject, error) {
	return s.adminHandler().AdminDeleteReview(ctx, request)
}

func (s *Server) AdminApproveReview(ctx context.Context, request openapi.AdminApproveReviewRequestObject) (openapi.AdminApproveReviewResponseObject, error) {
	return s.adminHandler().AdminApproveReview(ctx, request)
}

func (s *Server) AdminHideReview(ctx context.Context, request openapi.AdminHideReviewRequestObject) (openapi.AdminHideReviewResponseObject, error) {
	return s.adminHandler().AdminHideReview(ctx, request)
}

func (s *Server) AdminFeatureReview(ctx context.Context, request openapi.AdminFeatureReviewRequestObject) (openapi.AdminFeatureReviewResponseObject, error) {
	return s.adminHandler().AdminFeatureReview(ctx, request)
}

// Admin Users
func (s *Server) AdminListUsers(ctx context.Context, request openapi.AdminListUsersRequestObject) (openapi.AdminListUsersResponseObject, error) {
	return s.adminHandler().AdminListUsers(ctx, request)
}

func (s *Server) AdminBlockUser(ctx context.Context, request openapi.AdminBlockUserRequestObject) (openapi.AdminBlockUserResponseObject, error) {
	return s.adminHandler().AdminBlockUser(ctx, request)
}

// Admin Settings
func (s *Server) AdminGetSetting(ctx context.Context, request openapi.AdminGetSettingRequestObject) (openapi.AdminGetSettingResponseObject, error) {
	return s.adminHandler().AdminGetSetting(ctx, request)
}

func (s *Server) AdminUpsertSetting(ctx context.Context, request openapi.AdminUpsertSettingRequestObject) (openapi.AdminUpsertSettingResponseObject, error) {
	return s.adminHandler().AdminUpsertSetting(ctx, request)
}

func (s *Server) AdminGetStoreSettings(ctx context.Context, request openapi.AdminGetStoreSettingsRequestObject) (openapi.AdminGetStoreSettingsResponseObject, error) {
	return s.adminHandler().AdminGetStoreSettings(ctx, request)
}

func (s *Server) AdminUpdateStoreSettingsBrand(ctx context.Context, request openapi.AdminUpdateStoreSettingsBrandRequestObject) (openapi.AdminUpdateStoreSettingsBrandResponseObject, error) {
	return s.adminHandler().AdminUpdateStoreSettingsBrand(ctx, request)
}

func (s *Server) AdminUpdateStoreSettingsContact(ctx context.Context, request openapi.AdminUpdateStoreSettingsContactRequestObject) (openapi.AdminUpdateStoreSettingsContactResponseObject, error) {
	return s.adminHandler().AdminUpdateStoreSettingsContact(ctx, request)
}

func (s *Server) AdminUpdateStoreSettingsFooter(ctx context.Context, request openapi.AdminUpdateStoreSettingsFooterRequestObject) (openapi.AdminUpdateStoreSettingsFooterResponseObject, error) {
	return s.adminHandler().AdminUpdateStoreSettingsFooter(ctx, request)
}

func (s *Server) AdminUpdateStoreSettingsHomepage(ctx context.Context, request openapi.AdminUpdateStoreSettingsHomepageRequestObject) (openapi.AdminUpdateStoreSettingsHomepageResponseObject, error) {
	return s.adminHandler().AdminUpdateStoreSettingsHomepage(ctx, request)
}

func (s *Server) AdminUpdateStoreSettingsFloatingSupport(ctx context.Context, request openapi.AdminUpdateStoreSettingsFloatingSupportRequestObject) (openapi.AdminUpdateStoreSettingsFloatingSupportResponseObject, error) {
	return s.adminHandler().AdminUpdateStoreSettingsFloatingSupport(ctx, request)
}

// Admin Media
func (s *Server) AdminListMedia(ctx context.Context, request openapi.AdminListMediaRequestObject) (openapi.AdminListMediaResponseObject, error) {
	return s.adminHandler().AdminListMedia(ctx, request)
}

func (s *Server) AdminCreateMedia(ctx context.Context, request openapi.AdminCreateMediaRequestObject) (openapi.AdminCreateMediaResponseObject, error) {
	return s.adminHandler().AdminCreateMedia(ctx, request)
}

func (s *Server) AdminGetPresignedUrl(ctx context.Context, request openapi.AdminGetPresignedUrlRequestObject) (openapi.AdminGetPresignedUrlResponseObject, error) {
	return s.adminHandler().AdminGetPresignedUrl(ctx, request)
}

// Admin Banners & FAQs
func (s *Server) AdminListBanners(ctx context.Context, request openapi.AdminListBannersRequestObject) (openapi.AdminListBannersResponseObject, error) {
	return s.adminHandler().AdminListBanners(ctx, request)
}

func (s *Server) AdminSaveBanner(ctx context.Context, request openapi.AdminSaveBannerRequestObject) (openapi.AdminSaveBannerResponseObject, error) {
	return s.adminHandler().AdminSaveBanner(ctx, request)
}

func (s *Server) AdminListFaqs(ctx context.Context, request openapi.AdminListFaqsRequestObject) (openapi.AdminListFaqsResponseObject, error) {
	return s.adminHandler().AdminListFaqs(ctx, request)
}

func (s *Server) AdminSaveFaq(ctx context.Context, request openapi.AdminSaveFaqRequestObject) (openapi.AdminSaveFaqResponseObject, error) {
	return s.adminHandler().AdminSaveFaq(ctx, request)
}

// Admin Messages & Notifications
func (s *Server) AdminListMessages(ctx context.Context, request openapi.AdminListMessagesRequestObject) (openapi.AdminListMessagesResponseObject, error) {
	return s.adminHandler().AdminListMessages(ctx, request)
}

func (s *Server) AdminBroadcastMessage(ctx context.Context, request openapi.AdminBroadcastMessageRequestObject) (openapi.AdminBroadcastMessageResponseObject, error) {
	return s.adminHandler().AdminBroadcastMessage(ctx, request)
}

func (s *Server) AdminGetNotifications(ctx context.Context, request openapi.AdminGetNotificationsRequestObject) (openapi.AdminGetNotificationsResponseObject, error) {
	return s.adminHandler().AdminGetNotifications(ctx, request)
}

// Admin Products & Categories
func (s *Server) AdminListProducts(ctx context.Context, request openapi.AdminListProductsRequestObject) (openapi.AdminListProductsResponseObject, error) {
	return s.adminHandler().AdminListProducts(ctx, request)
}

func (s *Server) AdminListCategories(ctx context.Context, request openapi.AdminListCategoriesRequestObject) (openapi.AdminListCategoriesResponseObject, error) {
	return s.adminHandler().AdminListCategories(ctx, request)
}

// Auth Capability
func (s *Server) RegisterUser(ctx context.Context, request openapi.RegisterUserRequestObject) (openapi.RegisterUserResponseObject, error) {
	return s.authHandler().RegisterUser(ctx, request)
}

func (s *Server) LoginUser(ctx context.Context, request openapi.LoginUserRequestObject) (openapi.LoginUserResponseObject, error) {
	return s.authHandler().LoginUser(ctx, request)
}

func (s *Server) RefreshToken(ctx context.Context, request openapi.RefreshTokenRequestObject) (openapi.RefreshTokenResponseObject, error) {
	return s.authHandler().RefreshToken(ctx, request)
}

func (s *Server) LogoutUser(ctx context.Context, request openapi.LogoutUserRequestObject) (openapi.LogoutUserResponseObject, error) {
	return s.authHandler().LogoutUser(ctx, request)
}

func (s *Server) GetCurrentUser(ctx context.Context, request openapi.GetCurrentUserRequestObject) (openapi.GetCurrentUserResponseObject, error) {
	return s.authHandler().GetCurrentUser(ctx, request)
}

// Cart Capability
func (s *Server) GetMyCart(ctx context.Context, request openapi.GetMyCartRequestObject) (openapi.GetMyCartResponseObject, error) {
	return s.cartHandler().GetMyCart(ctx, request)
}

func (s *Server) AddToCart(ctx context.Context, request openapi.AddToCartRequestObject) (openapi.AddToCartResponseObject, error) {
	return s.cartHandler().AddToCart(ctx, request)
}

func (s *Server) UpdateCartItem(ctx context.Context, request openapi.UpdateCartItemRequestObject) (openapi.UpdateCartItemResponseObject, error) {
	return s.cartHandler().UpdateCartItem(ctx, request)
}

func (s *Server) RemoveCartItem(ctx context.Context, request openapi.RemoveCartItemRequestObject) (openapi.RemoveCartItemResponseObject, error) {
	return s.cartHandler().RemoveCartItem(ctx, request)
}

func (s *Server) ClearCart(ctx context.Context, request openapi.ClearCartRequestObject) (openapi.ClearCartResponseObject, error) {
	return s.cartHandler().ClearCart(ctx, request)
}

// Catalog Capability
func (s *Server) ListProducts(ctx context.Context, request openapi.ListProductsRequestObject) (openapi.ListProductsResponseObject, error) {
	return s.catalogHandler().ListProducts(ctx, request)
}

func (s *Server) CreateProduct(ctx context.Context, request openapi.CreateProductRequestObject) (openapi.CreateProductResponseObject, error) {
	return s.catalogHandler().CreateProduct(ctx, request)
}

func (s *Server) GetProductByID(ctx context.Context, request openapi.GetProductByIDRequestObject) (openapi.GetProductByIDResponseObject, error) {
	return s.catalogHandler().GetProductByID(ctx, request)
}

func (s *Server) UpdateProduct(ctx context.Context, request openapi.UpdateProductRequestObject) (openapi.UpdateProductResponseObject, error) {
	return s.catalogHandler().UpdateProduct(ctx, request)
}

func (s *Server) DeleteProduct(ctx context.Context, request openapi.DeleteProductRequestObject) (openapi.DeleteProductResponseObject, error) {
	return s.catalogHandler().DeleteProduct(ctx, request)
}

func (s *Server) ListCategories(ctx context.Context, request openapi.ListCategoriesRequestObject) (openapi.ListCategoriesResponseObject, error) {
	return s.catalogHandler().ListCategories(ctx, request)
}

func (s *Server) CreateCategory(ctx context.Context, request openapi.CreateCategoryRequestObject) (openapi.CreateCategoryResponseObject, error) {
	return s.catalogHandler().CreateCategory(ctx, request)
}

func (s *Server) ListTags(ctx context.Context, request openapi.ListTagsRequestObject) (openapi.ListTagsResponseObject, error) {
	return s.catalogHandler().ListTags(ctx, request)
}

func (s *Server) CreateTag(ctx context.Context, request openapi.CreateTagRequestObject) (openapi.CreateTagResponseObject, error) {
	return s.catalogHandler().CreateTag(ctx, request)
}

// Checkout Capability
func (s *Server) GetCheckoutPreview(ctx context.Context, request openapi.GetCheckoutPreviewRequestObject) (openapi.GetCheckoutPreviewResponseObject, error) {
	return s.checkoutHandler().GetCheckoutPreview(ctx, request)
}

func (s *Server) PreviewCheckout(ctx context.Context, request openapi.PreviewCheckoutRequestObject) (openapi.PreviewCheckoutResponseObject, error) {
	return s.checkoutHandler().PreviewCheckout(ctx, request)
}

func (s *Server) CreateCheckoutOrder(ctx context.Context, request openapi.CreateCheckoutOrderRequestObject) (openapi.CreateCheckoutOrderResponseObject, error) {
	return s.checkoutHandler().CreateCheckoutOrder(ctx, request)
}

func (s *Server) GetPaymentParams(ctx context.Context, request openapi.GetPaymentParamsRequestObject) (openapi.GetPaymentParamsResponseObject, error) {
	return s.checkoutHandler().GetPaymentParams(ctx, request)
}

func (s *Server) PostPaymentParams(ctx context.Context, request openapi.PostPaymentParamsRequestObject) (openapi.PostPaymentParamsResponseObject, error) {
	return s.checkoutHandler().PostPaymentParams(ctx, request)
}

func (s *Server) CreatePaymentOrder(ctx context.Context, request openapi.CreatePaymentOrderRequestObject) (openapi.CreatePaymentOrderResponseObject, error) {
	return s.checkoutHandler().CreatePaymentOrder(ctx, request)
}

func (s *Server) PaymentNotify(ctx context.Context, request openapi.PaymentNotifyRequestObject) (openapi.PaymentNotifyResponseObject, error) {
	return s.checkoutHandler().PaymentNotify(ctx, request)
}

func (s *Server) GetPaymentStatus(ctx context.Context, request openapi.GetPaymentStatusRequestObject) (openapi.GetPaymentStatusResponseObject, error) {
	return s.checkoutHandler().GetPaymentStatus(ctx, request)
}

// Content & Homepage Capability
func (s *Server) GetStaticPage(ctx context.Context, request openapi.GetStaticPageRequestObject) (openapi.GetStaticPageResponseObject, error) {
	return s.contentHandler().GetStaticPage(ctx, request)
}

func (s *Server) GetHomepageConfig(ctx context.Context, request openapi.GetHomepageConfigRequestObject) (openapi.GetHomepageConfigResponseObject, error) {
	return s.contentHandler().GetHomepageConfig(ctx, request)
}

// Importer Capability
func (s *Server) ExecuteImport(ctx context.Context, request openapi.ExecuteImportRequestObject) (openapi.ExecuteImportResponseObject, error) {
	return s.importerHandler().ExecuteImport(ctx, request)
}

// Lead Capability
func (s *Server) SubmitLead(ctx context.Context, request openapi.SubmitLeadRequestObject) (openapi.SubmitLeadResponseObject, error) {
	return s.leadHandler().SubmitLead(ctx, request)
}

func (s *Server) ListLeads(ctx context.Context, request openapi.ListLeadsRequestObject) (openapi.ListLeadsResponseObject, error) {
	return s.leadHandler().ListLeads(ctx, request)
}

// Media Capability
func (s *Server) UploadMedia(ctx context.Context, request openapi.UploadMediaRequestObject) (openapi.UploadMediaResponseObject, error) {
	return s.mediaHandler().UploadMedia(ctx, request)
}

func (s *Server) DeleteMedia(ctx context.Context, request openapi.DeleteMediaRequestObject) (openapi.DeleteMediaResponseObject, error) {
	return s.mediaHandler().DeleteMedia(ctx, request)
}

// Notification Capability
func (s *Server) ListNotifications(ctx context.Context, request openapi.ListNotificationsRequestObject) (openapi.ListNotificationsResponseObject, error) {
	return s.notificationHandler().ListNotifications(ctx, request)
}

func (s *Server) MarkAllNotificationsRead(ctx context.Context, request openapi.MarkAllNotificationsReadRequestObject) (openapi.MarkAllNotificationsReadResponseObject, error) {
	return s.notificationHandler().MarkAllNotificationsRead(ctx, request)
}

func (s *Server) GetUnreadNotificationCount(ctx context.Context, request openapi.GetUnreadNotificationCountRequestObject) (openapi.GetUnreadNotificationCountResponseObject, error) {
	return s.notificationHandler().GetUnreadNotificationCount(ctx, request)
}

func (s *Server) MarkNotificationRead(ctx context.Context, request openapi.MarkNotificationReadRequestObject) (openapi.MarkNotificationReadResponseObject, error) {
	return s.notificationHandler().MarkNotificationRead(ctx, request)
}

// Reviews Capability
func (s *Server) GetProductReviews(ctx context.Context, request openapi.GetProductReviewsRequestObject) (openapi.GetProductReviewsResponseObject, error) {
	return s.reviewsHandler().GetProductReviews(ctx, request)
}

func (s *Server) CreateReview(ctx context.Context, request openapi.CreateReviewRequestObject) (openapi.CreateReviewResponseObject, error) {
	return s.reviewsHandler().CreateReview(ctx, request)
}

func (s *Server) DeleteReview(ctx context.Context, request openapi.DeleteReviewRequestObject) (openapi.DeleteReviewResponseObject, error) {
	return s.reviewsHandler().DeleteReview(ctx, request)
}

// Orders Capability
func (s *Server) ListOrders(ctx context.Context, request openapi.ListOrdersRequestObject) (openapi.ListOrdersResponseObject, error) {
	return s.ordersHandler().ListOrders(ctx, request)
}

func (s *Server) GetOrderByID(ctx context.Context, request openapi.GetOrderByIDRequestObject) (openapi.GetOrderByIDResponseObject, error) {
	return s.ordersHandler().GetOrderByID(ctx, request)
}

func (s *Server) RequestOrderRefund(ctx context.Context, request openapi.RequestOrderRefundRequestObject) (openapi.RequestOrderRefundResponseObject, error) {
	return s.ordersHandler().RequestOrderRefund(ctx, request)
}

// Profile Capability
func (s *Server) GetAccountProfile(ctx context.Context, request openapi.GetAccountProfileRequestObject) (openapi.GetAccountProfileResponseObject, error) {
	return s.profileHandler().GetAccountProfile(ctx, request)
}

func (s *Server) UpdateAccountProfile(ctx context.Context, request openapi.UpdateAccountProfileRequestObject) (openapi.UpdateAccountProfileResponseObject, error) {
	return s.profileHandler().UpdateAccountProfile(ctx, request)
}

// User Capability
func (s *Server) GetMyProfile(ctx context.Context, request openapi.GetMyProfileRequestObject) (openapi.GetMyProfileResponseObject, error) {
	return s.userHandler().GetMyProfile(ctx, request)
}

func (s *Server) UpdateMyProfile(ctx context.Context, request openapi.UpdateMyProfileRequestObject) (openapi.UpdateMyProfileResponseObject, error) {
	return s.userHandler().UpdateMyProfile(ctx, request)
}

func (s *Server) ListUsers(ctx context.Context, request openapi.ListUsersRequestObject) (openapi.ListUsersResponseObject, error) {
	return s.userHandler().ListUsers(ctx, request)
}

func (s *Server) CreateAdminUser(ctx context.Context, request openapi.CreateAdminUserRequestObject) (openapi.CreateAdminUserResponseObject, error) {
	return s.userHandler().CreateAdminUser(ctx, request)
}

func (s *Server) GetUserByID(ctx context.Context, request openapi.GetUserByIDRequestObject) (openapi.GetUserByIDResponseObject, error) {
	return s.userHandler().GetUserByID(ctx, request)
}

func (s *Server) LockUser(ctx context.Context, request openapi.LockUserRequestObject) (openapi.LockUserResponseObject, error) {
	return s.userHandler().LockUser(ctx, request)
}

func (s *Server) UnlockUser(ctx context.Context, request openapi.UnlockUserRequestObject) (openapi.UnlockUserResponseObject, error) {
	return s.userHandler().UnlockUser(ctx, request)
}

// Wishlist Capability
func (s *Server) GetMyWishlist(ctx context.Context, request openapi.GetMyWishlistRequestObject) (openapi.GetMyWishlistResponseObject, error) {
	return s.wishlistHandler().GetMyWishlist(ctx, request)
}

func (s *Server) AddToWishlistDirect(ctx context.Context, request openapi.AddToWishlistDirectRequestObject) (openapi.AddToWishlistDirectResponseObject, error) {
	return s.wishlistHandler().AddToWishlistDirect(ctx, request)
}

func (s *Server) RemoveFromWishlistDirect(ctx context.Context, request openapi.RemoveFromWishlistDirectRequestObject) (openapi.RemoveFromWishlistDirectResponseObject, error) {
	return s.wishlistHandler().RemoveFromWishlistDirect(ctx, request)
}

func (s *Server) VoteWishlistItem(ctx context.Context, request openapi.VoteWishlistItemRequestObject) (openapi.VoteWishlistItemResponseObject, error) {
	return s.wishlistHandler().VoteWishlistItem(ctx, request)
}

func (s *Server) AddToWishlist(ctx context.Context, request openapi.AddToWishlistRequestObject) (openapi.AddToWishlistResponseObject, error) {
	return s.wishlistHandler().AddToWishlist(ctx, request)
}

func (s *Server) RemoveFromWishlist(ctx context.Context, request openapi.RemoveFromWishlistRequestObject) (openapi.RemoveFromWishlistResponseObject, error) {
	return s.wishlistHandler().RemoveFromWishlist(ctx, request)
}

// Additional Checkout methods
func (s *Server) CancelCheckoutOrder(ctx context.Context, request openapi.CancelCheckoutOrderRequestObject) (openapi.CancelCheckoutOrderResponseObject, error) {
	return s.checkoutHandler().CancelCheckoutOrder(ctx, request)
}

func (s *Server) GetCheckoutOrderStatus(ctx context.Context, request openapi.GetCheckoutOrderStatusRequestObject) (openapi.GetCheckoutOrderStatusResponseObject, error) {
	return s.checkoutHandler().GetCheckoutOrderStatus(ctx, request)
}

// Additional Content methods
func (s *Server) GetPublicHomepage(ctx context.Context, request openapi.GetPublicHomepageRequestObject) (openapi.GetPublicHomepageResponseObject, error) {
	return s.contentHandler().GetPublicHomepage(ctx, request)
}

func (s *Server) GetActiveFaqs(ctx context.Context, request openapi.GetActiveFaqsRequestObject) (openapi.GetActiveFaqsResponseObject, error) {
	return s.contentHandler().GetActiveFaqs(ctx, request)
}

func (s *Server) ListContentArticles(ctx context.Context, request openapi.ListContentArticlesRequestObject) (openapi.ListContentArticlesResponseObject, error) {
	return s.contentHandler().ListContentArticles(ctx, request)
}

func (s *Server) CreateContentArticle(ctx context.Context, request openapi.CreateContentArticleRequestObject) (openapi.CreateContentArticleResponseObject, error) {
	return s.contentHandler().CreateContentArticle(ctx, request)
}

func (s *Server) ListContentPages(ctx context.Context, request openapi.ListContentPagesRequestObject) (openapi.ListContentPagesResponseObject, error) {
	return s.contentHandler().ListContentPages(ctx, request)
}

func (s *Server) CreateContentPage(ctx context.Context, request openapi.CreateContentPageRequestObject) (openapi.CreateContentPageResponseObject, error) {
	return s.contentHandler().CreateContentPage(ctx, request)
}

// Additional Profile methods
func (s *Server) GetProfile(ctx context.Context, request openapi.GetProfileRequestObject) (openapi.GetProfileResponseObject, error) {
	return s.profileHandler().GetProfile(ctx, request)
}

func (s *Server) GetUserProfile(ctx context.Context, request openapi.GetUserProfileRequestObject) (openapi.GetUserProfileResponseObject, error) {
	return s.profileHandler().GetUserProfile(ctx, request)
}

func (s *Server) UpdateProfile(ctx context.Context, request openapi.UpdateProfileRequestObject) (openapi.UpdateProfileResponseObject, error) {
	return s.profileHandler().UpdateProfile(ctx, request)
}

func (s *Server) UpdateProfileEmail(ctx context.Context, request openapi.UpdateProfileEmailRequestObject) (openapi.UpdateProfileEmailResponseObject, error) {
	return s.profileHandler().UpdateProfileEmail(ctx, request)
}

func (s *Server) GetProfileSecurity(ctx context.Context, request openapi.GetProfileSecurityRequestObject) (openapi.GetProfileSecurityResponseObject, error) {
	return s.profileHandler().GetProfileSecurity(ctx, request)
}

func (s *Server) GetProfileSessions(ctx context.Context, request openapi.GetProfileSessionsRequestObject) (openapi.GetProfileSessionsResponseObject, error) {
	return s.profileHandler().GetProfileSessions(ctx, request)
}

func (s *Server) UpdateProfileNotifications(ctx context.Context, request openapi.UpdateProfileNotificationsRequestObject) (openapi.UpdateProfileNotificationsResponseObject, error) {
	return s.profileHandler().UpdateProfileNotifications(ctx, request)
}







// Catalog additional method
func (s *Server) ListCatalogProductModels(ctx context.Context, request openapi.ListCatalogProductModelsRequestObject) (openapi.ListCatalogProductModelsResponseObject, error) {
	return s.catalogHandler().ListCatalogProductModels(ctx, request)
}


// Temporary stubs for newly declared OpenAPI operations (to be wired in Batch 2 & Batch 3)
func (s *Server) AdminDeleteBanner(ctx context.Context, request openapi.AdminDeleteBannerRequestObject) (openapi.AdminDeleteBannerResponseObject, error) {
	return nil, fmt.Errorf("AdminDeleteBanner not implemented")
}

func (s *Server) AdminListAttributeDefinitions(ctx context.Context, request openapi.AdminListAttributeDefinitionsRequestObject) (openapi.AdminListAttributeDefinitionsResponseObject, error) {
	return s.catalogHandler().AdminListAttributeDefinitions(ctx, request)
}

func (s *Server) AdminCreateAttributeDefinition(ctx context.Context, request openapi.AdminCreateAttributeDefinitionRequestObject) (openapi.AdminCreateAttributeDefinitionResponseObject, error) {
	return s.catalogHandler().AdminCreateAttributeDefinition(ctx, request)
}

func (s *Server) AdminUpdateAttributeDefinition(ctx context.Context, request openapi.AdminUpdateAttributeDefinitionRequestObject) (openapi.AdminUpdateAttributeDefinitionResponseObject, error) {
	return s.catalogHandler().AdminUpdateAttributeDefinition(ctx, request)
}

func (s *Server) AdminDeactivateAttributeDefinition(ctx context.Context, request openapi.AdminDeactivateAttributeDefinitionRequestObject) (openapi.AdminDeactivateAttributeDefinitionResponseObject, error) {
	return s.catalogHandler().AdminDeactivateAttributeDefinition(ctx, request)
}

func (s *Server) AdminAddAttributeEnumValue(ctx context.Context, request openapi.AdminAddAttributeEnumValueRequestObject) (openapi.AdminAddAttributeEnumValueResponseObject, error) {
	return s.catalogHandler().AdminAddAttributeEnumValue(ctx, request)
}

func (s *Server) AdminDeactivateAttributeEnumValue(ctx context.Context, request openapi.AdminDeactivateAttributeEnumValueRequestObject) (openapi.AdminDeactivateAttributeEnumValueResponseObject, error) {
	return s.catalogHandler().AdminDeactivateAttributeEnumValue(ctx, request)
}

func (s *Server) AdminListCatalogCategories(ctx context.Context, request openapi.AdminListCatalogCategoriesRequestObject) (openapi.AdminListCatalogCategoriesResponseObject, error) {
	return s.catalogHandler().AdminListCatalogCategories(ctx, request)
}

func (s *Server) AdminCreateCatalogCategory(ctx context.Context, request openapi.AdminCreateCatalogCategoryRequestObject) (openapi.AdminCreateCatalogCategoryResponseObject, error) {
	return s.catalogHandler().AdminCreateCatalogCategory(ctx, request)
}

func (s *Server) AdminDeleteCatalogCategory(ctx context.Context, request openapi.AdminDeleteCatalogCategoryRequestObject) (openapi.AdminDeleteCatalogCategoryResponseObject, error) {
	return s.catalogHandler().AdminDeleteCatalogCategory(ctx, request)
}

func (s *Server) AdminUpdateCatalogCategory(ctx context.Context, request openapi.AdminUpdateCatalogCategoryRequestObject) (openapi.AdminUpdateCatalogCategoryResponseObject, error) {
	return s.catalogHandler().AdminUpdateCatalogCategory(ctx, request)
}

func (s *Server) AdminDeactivateCatalogCategory(ctx context.Context, request openapi.AdminDeactivateCatalogCategoryRequestObject) (openapi.AdminDeactivateCatalogCategoryResponseObject, error) {
	return s.catalogHandler().AdminDeactivateCatalogCategory(ctx, request)
}

func (s *Server) AdminListCatalogMasters(ctx context.Context, request openapi.AdminListCatalogMastersRequestObject) (openapi.AdminListCatalogMastersResponseObject, error) {
	return s.catalogHandler().AdminListCatalogMasters(ctx, request)
}

func (s *Server) AdminCreateCatalogMaster(ctx context.Context, request openapi.AdminCreateCatalogMasterRequestObject) (openapi.AdminCreateCatalogMasterResponseObject, error) {
	return s.catalogHandler().AdminCreateCatalogMaster(ctx, request)
}

func (s *Server) AdminUpdateCatalogMaster(ctx context.Context, request openapi.AdminUpdateCatalogMasterRequestObject) (openapi.AdminUpdateCatalogMasterResponseObject, error) {
	return s.catalogHandler().AdminUpdateCatalogMaster(ctx, request)
}

func (s *Server) AdminDeactivateCatalogMaster(ctx context.Context, request openapi.AdminDeactivateCatalogMasterRequestObject) (openapi.AdminDeactivateCatalogMasterResponseObject, error) {
	return s.catalogHandler().AdminDeactivateCatalogMaster(ctx, request)
}

func (s *Server) AdminListCatalogProductModels(ctx context.Context, request openapi.AdminListCatalogProductModelsRequestObject) (openapi.AdminListCatalogProductModelsResponseObject, error) {
	return s.catalogHandler().AdminListCatalogProductModels(ctx, request)
}

func (s *Server) AdminCreateCatalogProductModel(ctx context.Context, request openapi.AdminCreateCatalogProductModelRequestObject) (openapi.AdminCreateCatalogProductModelResponseObject, error) {
	return s.catalogHandler().AdminCreateCatalogProductModel(ctx, request)
}

func (s *Server) AdminDeleteCatalogProductModel(ctx context.Context, request openapi.AdminDeleteCatalogProductModelRequestObject) (openapi.AdminDeleteCatalogProductModelResponseObject, error) {
	return s.catalogHandler().AdminDeleteCatalogProductModel(ctx, request)
}

func (s *Server) AdminGetCatalogProductModel(ctx context.Context, request openapi.AdminGetCatalogProductModelRequestObject) (openapi.AdminGetCatalogProductModelResponseObject, error) {
	return s.catalogHandler().AdminGetCatalogProductModel(ctx, request)
}

func (s *Server) AdminUpdateCatalogProductModel(ctx context.Context, request openapi.AdminUpdateCatalogProductModelRequestObject) (openapi.AdminUpdateCatalogProductModelResponseObject, error) {
	return s.catalogHandler().AdminUpdateCatalogProductModel(ctx, request)
}

func (s *Server) AdminDiscontinueCatalogProductModel(ctx context.Context, request openapi.AdminDiscontinueCatalogProductModelRequestObject) (openapi.AdminDiscontinueCatalogProductModelResponseObject, error) {
	return s.catalogHandler().AdminDiscontinueCatalogProductModel(ctx, request)
}

func (s *Server) AdminUpdateCatalogProductModelMedia(ctx context.Context, request openapi.AdminUpdateCatalogProductModelMediaRequestObject) (openapi.AdminUpdateCatalogProductModelMediaResponseObject, error) {
	return s.catalogHandler().AdminUpdateCatalogProductModelMedia(ctx, request)
}

func (s *Server) AdminPublishCatalogProductModel(ctx context.Context, request openapi.AdminPublishCatalogProductModelRequestObject) (openapi.AdminPublishCatalogProductModelResponseObject, error) {
	return s.catalogHandler().AdminPublishCatalogProductModel(ctx, request)
}

func (s *Server) AdminUnpublishCatalogProductModel(ctx context.Context, request openapi.AdminUnpublishCatalogProductModelRequestObject) (openapi.AdminUnpublishCatalogProductModelResponseObject, error) {
	return s.catalogHandler().AdminUnpublishCatalogProductModel(ctx, request)
}

func (s *Server) AdminAddCatalogVariantDimension(ctx context.Context, request openapi.AdminAddCatalogVariantDimensionRequestObject) (openapi.AdminAddCatalogVariantDimensionResponseObject, error) {
	return s.catalogHandler().AdminAddCatalogVariantDimension(ctx, request)
}

func (s *Server) AdminUpdateCatalogVariantDimension(ctx context.Context, request openapi.AdminUpdateCatalogVariantDimensionRequestObject) (openapi.AdminUpdateCatalogVariantDimensionResponseObject, error) {
	return s.catalogHandler().AdminUpdateCatalogVariantDimension(ctx, request)
}

func (s *Server) AdminAddCatalogVariantDimensionValue(ctx context.Context, request openapi.AdminAddCatalogVariantDimensionValueRequestObject) (openapi.AdminAddCatalogVariantDimensionValueResponseObject, error) {
	return s.catalogHandler().AdminAddCatalogVariantDimensionValue(ctx, request)
}

func (s *Server) AdminDeactivateCatalogVariantDimensionValue(ctx context.Context, request openapi.AdminDeactivateCatalogVariantDimensionValueRequestObject) (openapi.AdminDeactivateCatalogVariantDimensionValueResponseObject, error) {
	return s.catalogHandler().AdminDeactivateCatalogVariantDimensionValue(ctx, request)
}

func (s *Server) AdminListCatalogModelVariants(ctx context.Context, request openapi.AdminListCatalogModelVariantsRequestObject) (openapi.AdminListCatalogModelVariantsResponseObject, error) {
	return s.catalogHandler().AdminListCatalogModelVariants(ctx, request)
}

func (s *Server) AdminCreateCatalogModelVariant(ctx context.Context, request openapi.AdminCreateCatalogModelVariantRequestObject) (openapi.AdminCreateCatalogModelVariantResponseObject, error) {
	return s.catalogHandler().AdminCreateCatalogModelVariant(ctx, request)
}

func (s *Server) AdminBulkUpdateVariantPrices(ctx context.Context, request openapi.AdminBulkUpdateVariantPricesRequestObject) (openapi.AdminBulkUpdateVariantPricesResponseObject, error) {
	return s.catalogHandler().AdminBulkUpdateVariantPrices(ctx, request)
}

func (s *Server) AdminGetCatalogVariant(ctx context.Context, request openapi.AdminGetCatalogVariantRequestObject) (openapi.AdminGetCatalogVariantResponseObject, error) {
	return s.catalogHandler().AdminGetCatalogVariant(ctx, request)
}

func (s *Server) AdminUpdateCatalogVariant(ctx context.Context, request openapi.AdminUpdateCatalogVariantRequestObject) (openapi.AdminUpdateCatalogVariantResponseObject, error) {
	return s.catalogHandler().AdminUpdateCatalogVariant(ctx, request)
}

func (s *Server) AdminActivateCatalogVariant(ctx context.Context, request openapi.AdminActivateCatalogVariantRequestObject) (openapi.AdminActivateCatalogVariantResponseObject, error) {
	return s.catalogHandler().AdminActivateCatalogVariant(ctx, request)
}

func (s *Server) AdminInactivateCatalogVariant(ctx context.Context, request openapi.AdminInactivateCatalogVariantRequestObject) (openapi.AdminInactivateCatalogVariantResponseObject, error) {
	return s.catalogHandler().AdminInactivateCatalogVariant(ctx, request)
}

func (s *Server) AdminUpdateCollect(ctx context.Context, request openapi.AdminUpdateCollectRequestObject) (openapi.AdminUpdateCollectResponseObject, error) {
	return s.adminHandler().AdminUpdateCollect(ctx, request)
}

func (s *Server) AdminDeleteFaq(ctx context.Context, request openapi.AdminDeleteFaqRequestObject) (openapi.AdminDeleteFaqResponseObject, error) {
	return nil, fmt.Errorf("AdminDeleteFaq not implemented")
}

func (s *Server) QueueAdminNotificationTest(ctx context.Context, request openapi.QueueAdminNotificationTestRequestObject) (openapi.QueueAdminNotificationTestResponseObject, error) {
	return s.notificationHandler().QueueAdminNotificationTest(ctx, request)
}

func (s *Server) AdminUpdateProductEditorial(ctx context.Context, request openapi.AdminUpdateProductEditorialRequestObject) (openapi.AdminUpdateProductEditorialResponseObject, error) {
	return nil, fmt.Errorf("AdminUpdateProductEditorial not implemented")
}

func (s *Server) AdminGetRefund(ctx context.Context, request openapi.AdminGetRefundRequestObject) (openapi.AdminGetRefundResponseObject, error) {
	return s.adminHandler().AdminGetRefund(ctx, request)
}

func (s *Server) AdminGetReview(ctx context.Context, request openapi.AdminGetReviewRequestObject) (openapi.AdminGetReviewResponseObject, error) {
	return s.adminHandler().AdminGetReview(ctx, request)
}

func (s *Server) GetCatalogAnnouncement(ctx context.Context, request openapi.GetCatalogAnnouncementRequestObject) (openapi.GetCatalogAnnouncementResponseObject, error) {
	return s.catalogHandler().GetCatalogAnnouncement(ctx, request)
}

func (s *Server) GetCatalogProductModelOptions(ctx context.Context, request openapi.GetCatalogProductModelOptionsRequestObject) (openapi.GetCatalogProductModelOptionsResponseObject, error) {
	return s.catalogHandler().GetCatalogProductModelOptions(ctx, request)
}

func (s *Server) ResolveCatalogProductModelVariant(ctx context.Context, request openapi.ResolveCatalogProductModelVariantRequestObject) (openapi.ResolveCatalogProductModelVariantResponseObject, error) {
	return s.catalogHandler().ResolveCatalogProductModelVariant(ctx, request)
}

func (s *Server) GetCatalogProductBuyMeta(ctx context.Context, request openapi.GetCatalogProductBuyMetaRequestObject) (openapi.GetCatalogProductBuyMetaResponseObject, error) {
	return s.catalogHandler().GetCatalogProductBuyMeta(ctx, request)
}

func (s *Server) SearchCatalog(ctx context.Context, request openapi.SearchCatalogRequestObject) (openapi.SearchCatalogResponseObject, error) {
	return s.catalogHandler().SearchCatalog(ctx, request)
}

func (s *Server) GetCatalogSettings(ctx context.Context, request openapi.GetCatalogSettingsRequestObject) (openapi.GetCatalogSettingsResponseObject, error) {
	return s.catalogHandler().GetCatalogSettings(ctx, request)
}

func (s *Server) UpdateContentArticle(ctx context.Context, request openapi.UpdateContentArticleRequestObject) (openapi.UpdateContentArticleResponseObject, error) {
	return s.contentHandler().UpdateContentArticle(ctx, request)
}

func (s *Server) ListPublicContentArticles(ctx context.Context, request openapi.ListPublicContentArticlesRequestObject) (openapi.ListPublicContentArticlesResponseObject, error) {
	return s.contentHandler().ListPublicContentArticles(ctx, request)
}


func (s *Server) ClearNotificationInbox(ctx context.Context, request openapi.ClearNotificationInboxRequestObject) (openapi.ClearNotificationInboxResponseObject, error) {
	return s.notificationHandler().ClearNotificationInbox(ctx, request)
}

// Content additions
func (s *Server) GetContentArticlePreview(ctx context.Context, request openapi.GetContentArticlePreviewRequestObject) (openapi.GetContentArticlePreviewResponseObject, error) {
	return s.contentHandler().GetContentArticlePreview(ctx, request)
}

func (s *Server) GetPublicContentArticle(ctx context.Context, request openapi.GetPublicContentArticleRequestObject) (openapi.GetPublicContentArticleResponseObject, error) {
	return s.contentHandler().GetPublicContentArticle(ctx, request)
}

func (s *Server) UpdateContentPage(ctx context.Context, request openapi.UpdateContentPageRequestObject) (openapi.UpdateContentPageResponseObject, error) {
	return s.contentHandler().UpdateContentPage(ctx, request)
}
