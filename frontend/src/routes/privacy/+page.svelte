<script>
  import Header from '$lib/components/Header.svelte';
  import Footer from '$lib/components/Footer.svelte';
</script>

<svelte:head>
    <title>Privacy Policy - hn30</title>
    <meta name="description" content="Learn how hn30 protects your privacy. We collect minimal data to provide our service and do not track you." />
    <meta property="og:title" content="Privacy Policy - hn30" />
    <meta property="og:description" content="Learn how hn30 protects your privacy. We collect minimal data to provide our service and do not track you." />
    <link rel="canonical" href="https://hn.yamanlabs.com/privacy" />
    <meta property="og:url" content="https://hn.yamanlabs.com/privacy" />
</svelte:head>

<div class="min-h-screen flex flex-col">
  <Header />

  <main class="flex-grow">
    <div class="max-w-3xl mx-auto px-4 py-8 sm:px-6 lg:px-8">
        <h1 class="text-4xl font-extrabold mb-2 text-[var(--color-primary-text)]">Privacy Policy</h1>
        <p class="text-sm text-[var(--color-secondary-text)] mb-8">Last updated: September 23, 2025</p>

        <div class="space-y-6 text-[var(--color-secondary-text)] leading-relaxed">
            <section>
                <h2 class="text-2xl font-bold mt-8 mb-4 text-[var(--color-primary-text)] border-b border-[var(--color-border)] pb-2">Our Commitment to Your Privacy</h2>
                <p>Welcome to hn30. We believe in privacy by design. This policy outlines the minimal information we collect and how we use it to provide our service while respecting your privacy.</p>
            </section>

            <section>
                <h2 class="text-2xl font-bold mt-8 mb-4 text-[var(--color-primary-text)] border-b border-[var(--color-border)] pb-2">Information We Collect</h2>
                <p>We only collect information that is essential for the functionality and security of our service.</p>
                
                <h3 class="text-xl font-semibold mt-6 mb-3 text-[var(--color-primary-text)]">1. Server Logs for Security</h3>
                <p>
                    When you use our service, our backend server temporarily logs technical information for security and abuse prevention. This includes:
                </p>
                <ul class="list-disc list-inside pl-4 space-y-1">
                    <li>Your <strong class="font-semibold text-[var(--color-primary-text)]">IP Address</strong></li>
                    <li>The endpoint requested (e.g., for summaries or images)</li>
                    <li>The time of the request</li>
                </ul>
                <p class="mt-2">
                    Your IP address is used for rate-limiting our AI summarization and image proxy features. This is crucial to prevent misuse, manage costs, and ensure the service remains available to everyone. These logs are not used to identify you and to create personal profiles.
                </p>

                <h3 class="text-xl font-semibold mt-6 mb-3 text-[var(--color-primary-text)]">2. Anonymous Website Analytics</h3>
                <p>
                    We use Plausible Analytics, a privacy-first analytics tool, to understand website traffic. Plausible does not use cookies or generate persistent identifiers. It collects only anonymous, aggregate data, including:
                </p>
                <ul class="list-disc list-inside pl-4 space-y-1">
                    <li>The pages visited</li>
                    <li>The source of traffic (e.g., another website)</li>
                    <li>General information like browser, OS, and country</li>
                </ul>

                <h3 class="text-xl font-semibold mt-6 mb-3 text-[var(--color-primary-text)]">3. Your Browser's Local Storage</h3>
                <p>
                    We use your browser's <code class="bg-[var(--color-background-dark-sections)] text-[var(--color-primary-text)] font-mono text-sm rounded px-2 py-1">localStorage</code> and <code class="bg-[var(--color-background-dark-sections)] text-[var(--color-primary-text)] font-mono text-sm rounded px-2 py-1">Indexed Database</code> to locally save your preferences, such as:
                </p>
                <ul class="list-disc list-inside pl-4 space-y-1">
                    <li>Your theme choice (light or dark mode)</li>
                    <li>Your bookmarked stories</li>
                </ul>
                <p class="mt-2">This data is stored only on your device and is never sent to our servers. Clearing your browser data will remove it.</p>
            </section>

            <section>
                <h2 class="text-2xl font-bold mt-8 mb-4 text-[var(--color-primary-text)] border-b border-[var(--color-border)] pb-2">How We Handle Third-Party Services</h2>
                <ul class="list-disc list-inside pl-4 space-y-2">
                    <li><strong class="font-semibold text-[var(--color-primary-text)]">AI Summaries (Google Gemini):</strong> When you request an AI summary, the content of the article is sent to the Google Gemini API. We do not send any personal information about you to Google.</li>
                    <li><strong class="font-semibold text-[var(--color-primary-text)]">Image Proxy:</strong> To protect your privacy from third-party sites, all article images are routed through our backend proxy. This prevents websites from seeing your IP address or using tracking cookies when images are loaded.</li>
                </ul>
            </section>

            <section>
                <h2 class="text-2xl font-bold mt-8 mb-4 text-[var(--color-primary-text)] border-b border-[var(--color-border)] pb-2">Data Security</h2>
                <p>We take reasonable measures to protect the minimal data we handle from loss, misuse, and unauthorized access.</p>
            </section>

            <section>
                <h2 class="text-2xl font-bold mt-8 mb-4 text-[var(--color-primary-text)] border-b border-[var(--color-border)] pb-2">Your Rights</h2>
                <p>You have the right to know what data we collect. If you have any questions about this policy, please don't hesitate to reach out.</p>
            </section>

            <section>
                <h2 class="text-2xl font-bold mt-8 mb-4 text-[var(--color-primary-text)] border-b border-[var(--color-border)] pb-2">Contact Us</h2>
                <p>If you have questions about this Privacy Policy, you can find our contact information on <a href="https://yaman.pro" target="_blank" rel="noopener noreferrer" class="text-[var(--color-primary-accent)] hover:underline font-semibold">yaman.pro</a>.</p>
            </section>
        </div>
    </div>
  </main>

  <Footer />
</div>
