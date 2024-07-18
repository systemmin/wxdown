class API {
    constructor() {
        this.baseURL = "/";
    }

    async request(endpoint, method = 'GET', body = null, headers = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const item = sessionStorage.getItem("Authorization");
        console.log('Authorization', item)
        const options = {
            method,
            headers: {
                'Authorization': item,
                ...headers
            }
        };

        if (body) {
            options.body = JSON.stringify(body);
        }

        try {
            const response = await fetch(url, options);
            if (response.status === 403) {
                alert("未授权！")
                setTimeout(() => {
                    window.location.href = "login.html";
                }, 500)
                return;
            }
            if (response.status === 401) {
                alert("账号或密码错误！")
                setTimeout(() => {
                    window.location.href = "login.html";
                }, 500)
                return;
            }
            if (!response.ok) {
                const errorData = await response.json();
                console.log(errorData)
                throw new Error(errorData.message || 'Something went wrong');
            }
            const data = await response.json();
            if (!data.success) {
                throw new Error(data.msg || 'Something went wrong');
            }
            return data
        } catch (error) {
            console.error('API request error:', error);
            throw error;
        }
    }

    get(endpoint) {
        return this.request(endpoint, 'GET', null);
    }

    post(endpoint, body) {
        return this.request(endpoint, 'POST', body, {"Content-Type": "application/json"});
    }

    put(endpoint, body, headers = {}) {
        return this.request(endpoint, 'PUT', body, headers);
    }

    delete(endpoint, headers = {}) {
        return this.request(endpoint, 'DELETE', null, headers);
    }
}


