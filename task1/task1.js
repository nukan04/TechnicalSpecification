fetch('https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1')
    .then(response => response.json())
    .then(data => {
        const tableBody = document.getElementById('cryptoData');

        data.forEach(currency => {
            const row = tableBody.insertRow();
            row.dataset.symbol = currency.symbol;
            row.innerHTML = `
                        <td>${currency.id}</td>
                        <td>${currency.symbol}</td>
                        <td>${currency.name}</td>
                    `;

            if (currency.symbol === 'usdt') {
                row.classList.add('usdt');
            }
        });
    })
    .catch(error => console.error('Error:', error));