<script lang="ts">
	import { onMount } from 'svelte';
	import { Button, Badge, FeatureCard } from '$lib/components';
	
	let scrollY = 0;
	
	onMount(() => {
		const handleScroll = () => {
			scrollY = window.scrollY;
		};
		window.addEventListener('scroll', handleScroll);
		return () => window.removeEventListener('scroll', handleScroll);
	});
</script>

<svelte:head>
	<title>LinkBio - Một Link, Vô Vàn Khả Năng | Nền Tảng Link in Bio Thông Minh</title>
	<meta name="description" content="Tạo trang bio link đẹp mắt, chuyên nghiệp trong 60 giây. Tùy chỉnh không giới hạn, phân tích chi tiết, và hoàn toàn miễn phí." />
</svelte:head>

<div class="min-h-screen bg-white">
	<!-- Navigation -->
	<nav class="fixed top-0 w-full z-50 transition-all duration-300" 
	     class:bg-white={scrollY > 50}
	     class:shadow-sm={scrollY > 50}
	     class:bg-transparent={scrollY <= 50}>
		<div class="container mx-auto px-6 py-4">
			<div class="flex justify-between items-center">
				<div class="flex items-center space-x-2">
					<div class="w-8 h-8 bg-gradient-to-br from-violet-600 to-indigo-600 rounded-lg flex items-center justify-center">
						<span class="text-white font-bold text-lg">L</span>
					</div>
					<span class="text-xl font-bold bg-gradient-to-r from-violet-600 to-indigo-600 bg-clip-text text-transparent">
						LinkBio
					</span>
				</div>
				<div class="flex items-center space-x-6">
					<a href="/auth/login" class="text-gray-600 hover:text-gray-900 font-medium transition-colors">
						Đăng nhập
					</a>
					<a href="/auth/register" 
					   class="bg-gradient-to-r from-violet-600 to-indigo-600 text-white px-6 py-2.5 rounded-full font-medium hover:shadow-lg hover:scale-105 transition-all duration-200">
						Bắt đầu miễn phí
					</a>
				</div>
			</div>
		</div>
	</nav>

	<!-- Hero Section -->
	<section class="pt-32 pb-20 px-6 relative overflow-hidden">
		<!-- Background Elements -->
		<div class="absolute inset-0 -z-10">
			<div class="absolute top-20 left-10 w-72 h-72 bg-violet-200 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob"></div>
			<div class="absolute top-40 right-10 w-72 h-72 bg-indigo-200 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-2000"></div>
			<div class="absolute -bottom-8 left-1/2 w-72 h-72 bg-pink-200 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-4000"></div>
		</div>

		<div class="container mx-auto max-w-6xl">
			<div class="text-center space-y-8">
				<!-- Badge -->
				<Badge variant="default" class="inline-flex items-center space-x-2">
					<span class="w-2 h-2 bg-violet-600 rounded-full animate-pulse"></span>
					<span>Được tin dùng bởi 10,000+ creators</span>
				</Badge>

				<!-- Main Heading -->
				<h1 class="text-6xl md:text-7xl font-bold leading-tight">
					<span class="bg-gradient-to-r from-violet-600 via-indigo-600 to-purple-600 bg-clip-text text-transparent">
						Một Link.
					</span>
					<br />
					<span class="text-gray-900">Vô Vàn Khả Năng.</span>
				</h1>

				<!-- Subheading -->
				<p class="text-xl md:text-2xl text-gray-600 max-w-3xl mx-auto leading-relaxed">
					Tạo trang bio link đẹp mắt trong <span class="font-semibold text-violet-600">60 giây</span>. 
					Không cần code, không giới hạn links, và hoàn toàn <span class="font-semibold text-violet-600">miễn phí</span>.
				</p>

				<!-- CTA Buttons -->
				<div class="flex flex-col sm:flex-row items-center justify-center gap-4 pt-4">
					<a href="/auth/register" 
					   class="group inline-flex items-center space-x-2 bg-gradient-to-r from-violet-600 to-indigo-600 text-white px-8 py-4 rounded-full text-lg font-semibold hover:shadow-2xl hover:scale-105 transition-all duration-200">
						<span>Tạo Link Miễn Phí</span>
						<svg class="w-5 h-5 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"></path>
						</svg>
					</a>
					<a href="/auth/register" 
					   class="inline-flex items-center text-gray-700 px-8 py-4 rounded-full text-lg font-semibold border-2 border-gray-300 hover:border-violet-600 hover:text-violet-600 transition-all duration-200">
						Xem Demo
					</a>
				</div>

				<!-- Trust Indicators -->
				<div class="flex items-center justify-center space-x-8 pt-8 text-sm text-gray-500">
					<div class="flex items-center space-x-2">
						<svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
						</svg>
						<span>Không cần thẻ tín dụng</span>
					</div>
					<div class="flex items-center space-x-2">
						<svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
						</svg>
						<span>Setup trong 60 giây</span>
					</div>
					<div class="flex items-center space-x-2">
						<svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
						</svg>
						<span>Miễn phí mãi mãi</span>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Features Section -->
	<section class="py-20 px-6 bg-gray-50">
		<div class="container mx-auto max-w-6xl">
			<div class="text-center mb-16">
				<h2 class="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
					Tại sao chọn <span class="bg-gradient-to-r from-violet-600 to-indigo-600 bg-clip-text text-transparent">LinkBio</span>?
				</h2>
				<p class="text-xl text-gray-600">Vượt trội hơn đối thủ với những tính năng độc quyền</p>
			</div>

			<div class="grid md:grid-cols-3 gap-8">
				<!-- Feature 1 -->
				<div class="bg-white p-8 rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300 group">
					<div class="w-14 h-14 bg-gradient-to-br from-violet-500 to-indigo-500 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
						<svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
						</svg>
					</div>
					<h3 class="text-xl font-bold text-gray-900 mb-3">Siêu Nhanh</h3>
					<p class="text-gray-600 leading-relaxed">
						Trang của bạn tải trong <span class="font-semibold text-violet-600">dưới 1 giây</span>. 
						Tối ưu SEO và mobile-first để tiếp cận tối đa khách hàng.
					</p>
				</div>

				<!-- Feature 2 -->
				<div class="bg-white p-8 rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300 group">
					<div class="w-14 h-14 bg-gradient-to-br from-pink-500 to-rose-500 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
						<svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"></path>
						</svg>
					</div>
					<h3 class="text-xl font-bold text-gray-900 mb-3">Tùy Chỉnh Không Giới Hạn</h3>
					<p class="text-gray-600 leading-relaxed">
						Thay đổi màu sắc, font chữ, layout theo ý muốn. 
						<span class="font-semibold text-pink-600">20+ themes</span> đẹp mắt sẵn có.
					</p>
				</div>

				<!-- Feature 3 -->
				<div class="bg-white p-8 rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300 group">
					<div class="w-14 h-14 bg-gradient-to-br from-amber-500 to-orange-500 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
						<svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
						</svg>
					</div>
					<h3 class="text-xl font-bold text-gray-900 mb-3">Phân Tích Chi Tiết</h3>
					<p class="text-gray-600 leading-relaxed">
						Theo dõi clicks, views, và engagement real-time. 
						Hiểu rõ khách hàng để <span class="font-semibold text-amber-600">tăng conversion</span>.
					</p>
				</div>

				<!-- Feature 4 -->
				<div class="bg-white p-8 rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300 group">
					<div class="w-14 h-14 bg-gradient-to-br from-emerald-500 to-teal-500 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
						<svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
						</svg>
					</div>
					<h3 class="text-xl font-bold text-gray-900 mb-3">Bảo Mật Tuyệt Đối</h3>
					<p class="text-gray-600 leading-relaxed">
						SSL encryption, GDPR compliant. 
						Dữ liệu của bạn được <span class="font-semibold text-emerald-600">bảo vệ 100%</span>.
					</p>
				</div>

				<!-- Feature 5 -->
				<div class="bg-white p-8 rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300 group">
					<div class="w-14 h-14 bg-gradient-to-br from-blue-500 to-cyan-500 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
						<svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
						</svg>
					</div>
					<h3 class="text-xl font-bold text-gray-900 mb-3">SEO Tối Ưu</h3>
					<p class="text-gray-600 leading-relaxed">
						Tự động tối ưu cho Google. 
						Trang của bạn <span class="font-semibold text-blue-600">dễ tìm kiếm</span> hơn đối thủ.
					</p>
				</div>

				<!-- Feature 6 -->
				<div class="bg-white p-8 rounded-2xl shadow-sm hover:shadow-xl transition-all duration-300 group">
					<div class="w-14 h-14 bg-gradient-to-br from-purple-500 to-fuchsia-500 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
						<svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7"></path>
						</svg>
					</div>
					<h3 class="text-xl font-bold text-gray-900 mb-3">Miễn Phí Mãi Mãi</h3>
					<p class="text-gray-600 leading-relaxed">
						Không giới hạn links, không giới hạn clicks. 
						<span class="font-semibold text-purple-600">100% miễn phí</span>, không ẩn phí.
					</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Comparison Section -->
	<section class="py-20 px-6">
		<div class="container mx-auto max-w-5xl">
			<div class="text-center mb-16">
				<h2 class="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
					So sánh với <span class="text-gray-400 line-through">đối thủ</span>
				</h2>
				<p class="text-xl text-gray-600">Xem sự khác biệt rõ ràng</p>
			</div>

			<div class="bg-white rounded-2xl shadow-xl overflow-hidden">
				<div class="grid grid-cols-3 gap-px bg-gray-200">
					<!-- Header -->
					<div class="bg-white p-6"></div>
					<div class="bg-gradient-to-br from-violet-50 to-indigo-50 p-6 text-center">
						<div class="font-bold text-xl text-violet-600">LinkBio</div>
						<div class="text-sm text-gray-600 mt-1">Của chúng tôi</div>
					</div>
					<div class="bg-white p-6 text-center">
						<div class="font-bold text-xl text-gray-400">Đối thủ</div>
						<div class="text-sm text-gray-500 mt-1">Linktree, Bio.fm...</div>
					</div>

					<!-- Row 1 -->
					<div class="bg-white p-6 font-medium text-gray-700">Giá</div>
					<div class="bg-gradient-to-br from-violet-50 to-indigo-50 p-6 text-center">
						<span class="text-2xl font-bold text-violet-600">$0</span>
						<div class="text-sm text-gray-600">Miễn phí mãi mãi</div>
					</div>
					<div class="bg-white p-6 text-center">
						<span class="text-2xl font-bold text-gray-400">$5-29</span>
						<div class="text-sm text-gray-500">Mỗi tháng</div>
					</div>

					<!-- Row 2 -->
					<div class="bg-white p-6 font-medium text-gray-700">Số lượng links</div>
					<div class="bg-gradient-to-br from-violet-50 to-indigo-50 p-6 text-center">
						<svg class="w-8 h-8 text-green-500 mx-auto" fill="currentColor" viewBox="0 0 20 20">
							<path d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"></path>
						</svg>
						<div class="text-sm text-gray-600 mt-2">Không giới hạn</div>
					</div>
					<div class="bg-white p-6 text-center">
						<svg class="w-8 h-8 text-red-400 mx-auto" fill="currentColor" viewBox="0 0 20 20">
							<path d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"></path>
						</svg>
						<div class="text-sm text-gray-500 mt-2">Giới hạn 5-10</div>
					</div>

					<!-- Row 3 -->
					<div class="bg-white p-6 font-medium text-gray-700">Tùy chỉnh theme</div>
					<div class="bg-gradient-to-br from-violet-50 to-indigo-50 p-6 text-center">
						<svg class="w-8 h-8 text-green-500 mx-auto" fill="currentColor" viewBox="0 0 20 20">
							<path d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"></path>
						</svg>
						<div class="text-sm text-gray-600 mt-2">20+ themes</div>
					</div>
					<div class="bg-white p-6 text-center">
						<svg class="w-8 h-8 text-red-400 mx-auto" fill="currentColor" viewBox="0 0 20 20">
							<path d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"></path>
						</svg>
						<div class="text-sm text-gray-500 mt-2">Cần trả phí</div>
					</div>

					<!-- Row 4 -->
					<div class="bg-white p-6 font-medium text-gray-700">Analytics</div>
					<div class="bg-gradient-to-br from-violet-50 to-indigo-50 p-6 text-center">
						<svg class="w-8 h-8 text-green-500 mx-auto" fill="currentColor" viewBox="0 0 20 20">
							<path d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"></path>
						</svg>
						<div class="text-sm text-gray-600 mt-2">Chi tiết real-time</div>
					</div>
					<div class="bg-white p-6 text-center">
						<svg class="w-8 h-8 text-yellow-400 mx-auto" fill="currentColor" viewBox="0 0 20 20">
							<path d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V7z"></path>
						</svg>
						<div class="text-sm text-gray-500 mt-2">Cơ bản</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- CTA Section -->
	<section class="py-20 px-6 bg-gradient-to-br from-violet-600 via-indigo-600 to-purple-600 relative overflow-hidden">
		<!-- Decorative elements -->
		<div class="absolute inset-0 opacity-10">
			<div class="absolute top-0 left-0 w-96 h-96 bg-white rounded-full filter blur-3xl"></div>
			<div class="absolute bottom-0 right-0 w-96 h-96 bg-white rounded-full filter blur-3xl"></div>
		</div>

		<div class="container mx-auto max-w-4xl text-center relative z-10">
			<h2 class="text-4xl md:text-5xl font-bold text-white mb-6">
				Sẵn sàng tạo link của bạn?
			</h2>
			<p class="text-xl text-violet-100 mb-10 max-w-2xl mx-auto">
				Tham gia cùng 10,000+ creators đang sử dụng LinkBio để phát triển thương hiệu cá nhân
			</p>
			<div class="flex flex-col sm:flex-row items-center justify-center gap-4">
				<a href="/auth/register" 
				   class="bg-white text-violet-600 px-10 py-4 rounded-full text-lg font-bold hover:bg-gray-50 hover:scale-105 transition-all duration-200 shadow-xl">
					Tạo Link Miễn Phí Ngay
				</a>
				<a href="#demo" 
				   class="text-white px-10 py-4 rounded-full text-lg font-semibold border-2 border-white hover:bg-white hover:text-violet-600 transition-all duration-200">
					Xem Demo
				</a>
			</div>
			<p class="text-violet-100 mt-6 text-sm">
				✨ Không cần thẻ tín dụng • Setup trong 60 giây • Miễn phí mãi mãi
			</p>
		</div>
	</section>

	<!-- Footer -->
	<footer class="bg-gray-900 text-gray-400 py-12 px-6">
		<div class="container mx-auto max-w-6xl">
			<div class="flex flex-col md:flex-row justify-between items-center">
				<div class="mb-6 md:mb-0">
					<div class="flex items-center space-x-2 mb-2">
						<div class="w-8 h-8 bg-gradient-to-br from-violet-600 to-indigo-600 rounded-lg flex items-center justify-center">
							<span class="text-white font-bold text-lg">L</span>
						</div>
						<span class="text-xl font-bold text-white">LinkBio</span>
					</div>
					<p class="text-sm">© 2024 LinkBio. All rights reserved.</p>
				</div>
				<div class="flex space-x-8">
					<a href="/privacy" class="hover:text-white transition-colors">Privacy</a>
					<a href="/terms" class="hover:text-white transition-colors">Terms</a>
					<a href="/contact" class="hover:text-white transition-colors">Contact</a>
				</div>
			</div>
		</div>
	</footer>
</div>

<style>
	@keyframes blob {
		0% {
			transform: translate(0px, 0px) scale(1);
		}
		33% {
			transform: translate(30px, -50px) scale(1.1);
		}
		66% {
			transform: translate(-20px, 20px) scale(0.9);
		}
		100% {
			transform: translate(0px, 0px) scale(1);
		}
	}

	.animate-blob {
		animation: blob 7s infinite;
	}

	.animation-delay-2000 {
		animation-delay: 2s;
	}

	.animation-delay-4000 {
		animation-delay: 4s;
	}
</style>
