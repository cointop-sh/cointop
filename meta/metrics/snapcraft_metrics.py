import configparser
import pymacaroons
import requests
import json
from datetime import date, timedelta

# Generate config with "package_metrics" ACL:
# $ snapcraft export-login --snaps cointop --acls package_metrics --expires '2020-12-01 00:00:00' snapcraft.cfg
config_path = 'snapcraft.cfg'
days = 30

today = date.today()
start = (today - timedelta(days=days)).strftime("%Y-%m-%d")
end = today.strftime("%Y-%m-%d")

config = configparser.ConfigParser()
config.read(config_path)

macaroon = config['login.ubuntu.com']['macaroon']
unbound_discharge = config['login.ubuntu.com']['unbound_discharge']

root_macaroon = pymacaroons.Macaroon.deserialize(macaroon)
unbound = pymacaroons.Macaroon.deserialize(unbound_discharge)
bound = root_macaroon.prepare_for_request(unbound)
discharge_macaroon_raw = bound.serialize()
auth = "Macaroon root={}, discharge={}".format(macaroon, discharge_macaroon_raw)

# https://dashboard.snapcraft.io/docs/api/snap.html#obtaining-information-about-a-snap
url = 'https://dashboard.snapcraft.io/dev/api/snaps/info/cointop'
r = requests.get(url, headers={'Authorization': auth, 'Content-Type': 'application/json'})
data = r.json()
snap_id = data['snap_id']

# https://dashboard.snapcraft.io/docs/api/snap.html#post--dev-api-snaps-metrics
url = 'https://dashboard.snapcraft.io/dev/api/snaps/metrics'
payload = {
    'filters': [
        {
            'snap_id': snap_id,
            'metric_name': 'installed_base_by_channel',
            'start': start,
            'end': end
        }
    ]
}
r = requests.post(
        url,
        data=json.dumps(payload),
        headers={
            'Authorization': auth,
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        }
    )
data = r.json()
series = data['metrics'][0]['series']
values = [x for x in series if x['name'] == 'stable'][0]['values']
total = sum(values)

print(total)
