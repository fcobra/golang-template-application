const fetchTemplate = async (name) => {
    const response = await fetch(`components/${name}.html`);
    if (!response.ok) {
        throw new Error(`Could not load template for ${name}`);
    }
    return await response.text();
};

async function init() {
    try {
        const [homeTemplate, loginTemplate, catalogTemplate] = await Promise.all([
            fetchTemplate('home-page'),
            fetchTemplate('login-page'),
            fetchTemplate('catalog-page'),
        ]);

        const App = {
            data() {
                return {
                    user: null,
                    currentPage: 'home',
                    isReady: true, // App is ready to be rendered
                    auth: { email: 'test@example.com', password: 'password123' },
                    data: { key: 'test-key', value: 'test-value' },
                    catalog: { items: [], loading: false, error: null },
                    message: '',
                    messageType: '',
                };
            },
            computed: {
                currentComponent() {
                    if (this.user) {
                        return this.currentPage === 'catalog' ? 'catalog-page' : 'home-page';
                    }
                    return 'login-page';
                }
            },
            mounted() {
                window.addEventListener('hashchange', this.handleRouteChange);
                this.checkAuthStatus();
            },
            methods: {
                handleRouteChange() {
                    const route = window.location.hash.slice(1);
                    if (route.startsWith('/catalog')) {
                        this.currentPage = 'catalog';
                        if (this.user) {
                            this.getCatalog();
                        }
                    } else {
                        this.currentPage = 'home';
                    }
                },
                async checkAuthStatus() {
                    try {
                        const response = await fetch('http://localhost:8080/api/v1/auth/me', { credentials: 'include' });
                        this.user = response.ok ? await response.json() : null;
                    } catch (error) {
                        this.user = null;
                    } finally {
                        this.handleRouteChange();
                    }
                },
                async login() {
                    this.clearMessage();
                    try {
                        const response = await fetch('http://localhost:8080/api/v1/auth/login', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify(this.auth),
                            credentials: 'include',
                        });
                        if (!response.ok) {
                            throw new Error(`Login failed: ${response.status} ${await response.text()}`);
                        }
                        this.user = await response.json();
                        this.showMessage('Login successful!', 'success');
                        window.location.hash = this.currentPage === 'catalog' ? '/catalog' : '/';
                    } catch (error) {
                        this.showMessage(error.message, 'error');
                    }
                },
                async logout() {
                    this.clearMessage();
                    try {
                        await fetch('http://localhost:8080/api/v1/auth/logout', { method: 'POST', credentials: 'include' });
                        this.user = null;
                        window.location.hash = '/';
                        this.showMessage('You have been logged out.', 'success');
                    } catch (error) {
                        this.showMessage(error.message, 'error');
                    }
                },
                async getCatalog() {
                    if (this.catalog.items.length > 0) return;
                    this.catalog.loading = true;
                    this.catalog.error = null;
                    try {
                        const response = await fetch('http://localhost:8080/api/v1/catalog', { credentials: 'include' });
                        if (response.status === 401) throw new Error('You must be logged in to view the catalog.');
                        if (!response.ok) throw new Error('Failed to fetch catalog data.');
                        this.catalog.items = await response.json();
                    } catch (error) {
                        this.catalog.error = error.message;
                    } finally {
                        this.catalog.loading = false;
                    }
                },
                async postData() {
                    this.clearMessage();
                    try {
                        const response = await fetch('http://localhost:8080/api/v1/data', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify(this.data),
                            credentials: 'include',
                        });
                        if (!response.ok) {
                            throw new Error(`Post data failed: ${response.status} ${await response.text()}`);
                        }
                        this.showMessage('Data posted successfully!', 'success');
                    } catch (error) {
                        this.showMessage(error.message, 'error');
                    }
                },
                showMessage(msg, type) {
                    this.message = msg;
                    this.messageType = type;
                    setTimeout(() => this.clearMessage(), 5000);
                },
                clearMessage() {
                    this.message = '';
                    this.messageType = '';
                }
            }
        };

        const app = Vue.createApp(App);

        // Register components with the fetched templates
        app.component('home-page', { props: ['appState'], template: homeTemplate });
        app.component('login-page', { props: ['appState'], template: loginTemplate });
        app.component('catalog-page', { props: ['appState'], template: catalogTemplate });

        app.mount('#app');

    } catch (error) {
        document.getElementById('app').innerHTML = `<div class="message error">Fatal Error: Could not load UI. Please check the console.</div>`;
        console.error(error);
    }
}

// Start the initialization process
init();
