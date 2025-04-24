import axios from 'axios'

const getSettings = async () => {
  let response = await axios.get(`/filebrowser/db/settings?size=1000`)
  let settings = {}
  for(let x of response.results) {
    if(x.setting_value === 'true' || x.setting_value === 'false')
      x.setting_value = JSON.parse(x.setting_value)
    settings[x.setting_key] = x.setting_value
  }
  return settings
}

export {
  getSettings
}