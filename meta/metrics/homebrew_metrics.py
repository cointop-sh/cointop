import requests

days = 30

# https://formulae.brew.sh/api/formula/cointop.json
url = f'https://formulae.brew.sh/api/analytics/install/{days}d.json'

r = requests.get(url)
data = r.json()
entry = [x for x in data['items'] if x['formula'] == 'cointop'][0]
total = entry['count']

print(total)
