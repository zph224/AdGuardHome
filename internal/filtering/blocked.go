package filtering

import (
	"encoding/json"
	"net/http"

	"github.com/AdguardTeam/AdGuardHome/internal/aghhttp"
	"github.com/AdguardTeam/golibs/log"
	"github.com/AdguardTeam/urlfilter/rules"
)

var serviceRules map[string][]*rules.NetworkRule // service name -> filtering rules

type svc struct {
	name  string
	rules []string
}

// Keep in sync with:
// client/src/helpers/constants.js
// client/src/components/ui/Icons.js
var serviceRulesArray = []svc{
	{"whatsapp", []string{"||whatsapp.net^", "||whatsapp.com^"}},
	{"facebook", []string{
		"||facebook.com^",
		"||facebook.net^",
		"||fbcdn.net^",
		"||accountkit.com^",
		"||fb.me^",
		"||fb.com^",
		"||fbsbx.com^",
		"||messenger.com^",
		"||facebookcorewwwi.onion^",
		"||fbcdn.com^",
		"||fb.watch^",
	}},
	{"twitter", []string{"||twitter.com^", "||twttr.com^", "||t.co^", "||twimg.com^"}},
	{"youtube", []string{
		"||youtube.com^",
		"||ytimg.com^",
		"||youtu.be^",
		"||googlevideo.com^",
		"||youtubei.googleapis.com^",
		"||youtube-nocookie.com^",
		"||youtube",
	}},
	{"twitch", []string{"||twitch.tv^", "||ttvnw.net^", "||jtvnw.net^", "||twitchcdn.net^"}},
	{"netflix", []string{"||nflxext.com^", "||netflix.com^", "||nflximg.net^", "||nflxvideo.net^", "||nflxso.net^"}},
	{"instagram", []string{"||instagram.com^", "||cdninstagram.com^"}},
	{"snapchat", []string{
		"||snapchat.com^",
		"||sc-cdn.net^",
		"||snap-dev.net^",
		"||snapkit.co",
		"||snapads.com^",
		"||impala-media-production.s3.amazonaws.com^",
	}},
	{"discord", []string{"||discord.gg^", "||discordapp.net^", "||discordapp.com^", "||discord.com^", "||discord.media^"}},
	{"ok", []string{"||ok.ru^"}},
	{"skype", []string{"||skype.com^", "||skypeassets.com^"}},
	{"vk", []string{"||vk.com^", "||userapi.com^", "||vk-cdn.net^", "||vkuservideo.net^"}},
	{"origin", []string{"||origin.com^", "||signin.ea.com^", "||accounts.ea.com^"}},
	{"steam", []string{
		"||steam.com^",
		"||steampowered.com^",
		"||steamcommunity.com^",
		"||steamstatic.com^",
		"||steamstore-a.akamaihd.net^",
		"||steamcdn-a.akamaihd.net^",
	}},
	{"epic_games", []string{"||epicgames.com^", "||easyanticheat.net^", "||easy.ac^", "||eac-cdn.com^"}},
	{"reddit", []string{"||reddit.com^", "||redditstatic.com^", "||redditmedia.com^", "||redd.it^"}},
	{"mail_ru", []string{"||mail.ru^"}},
	{"cloudflare", []string{
		"||cloudflare.com^",
		"||cloudflare-dns.com^",
		"||cloudflare.net^",
		"||cloudflareinsights.com^",
		"||cloudflarestream.com^",
		"||cloudflareresolve.com^",
		"||cloudflareclient.com^",
		"||cloudflarebolt.com^",
		"||cloudflarestatus.com^",
		"||cloudflare.cn^",
		"||one.one^",
		"||warp.plus^",
		"||1.1.1.1^",
		"||dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion^",
	}},
	{"amazon", []string{
		"||amazon.com^",
		"||media-amazon.com^",
		"||primevideo.com^",
		"||amazontrust.com^",
		"||images-amazon.com^",
		"||ssl-images-amazon.com^",
		"||amazonpay.com^",
		"||amazonpay.in^",
		"||amazon-adsystem.com^",
		"||a2z.com^",
		"||amazon.ae^",
		"||amazon.ca^",
		"||amazon.cn^",
		"||amazon.de^",
		"||amazon.es^",
		"||amazon.fr^",
		"||amazon.in^",
		"||amazon.it^",
		"||amazon.nl^",
		"||amazon.com.au^",
		"||amazon.com.br^",
		"||amazon.co.jp^",
		"||amazon.com.mx^",
		"||amazon.co.uk^",
		"||createspace.com^",
		"||aws",
	}},
	{"ebay", []string{
		"||ebay.com^",
		"||ebayimg.com^",
		"||ebaystatic.com^",
		"||ebaycdn.net^",
		"||ebayinc.com^",
		"||ebay.at^",
		"||ebay.be^",
		"||ebay.ca^",
		"||ebay.ch^",
		"||ebay.cn^",
		"||ebay.de^",
		"||ebay.es^",
		"||ebay.fr^",
		"||ebay.ie^",
		"||ebay.in^",
		"||ebay.it^",
		"||ebay.ph^",
		"||ebay.pl^",
		"||ebay.nl^",
		"||ebay.com.au^",
		"||ebay.com.cn^",
		"||ebay.com.hk^",
		"||ebay.com.my^",
		"||ebay.com.sg^",
		"||ebay.co.uk^",
	}},
	{"tiktok", []string{
		"||tiktok.com^",
		"||tiktokcdn.com^",
		"||musical.ly^",
		"||snssdk.com^",
		"||amemv.com^",
		"||toutiao.com^",
		"||ixigua.com^",
		"||pstatp.com^",
		"||ixiguavideo.com^",
		"||toutiaocloud.com^",
		"||toutiaocloud.net^",
		"||bdurl.com^",
		"||bytecdn.cn^",
		"||byteimg.com^",
		"||ixigua.com^",
		"||muscdn.com^",
		"||bytedance.map.fastly.net^",
		"||douyin.com^",
		"||tiktokv.com^",
	}},
	{"vimeo", []string{
		"||vimeo.com^",
		"||vimeocdn.com^",
		"*vod-adaptive.akamaized.net^",
	}},
	{"pinterest", []string{
		"||pinterest.*^",
		"||pinimg.com^",
	}},
	{"imgur", []string{
		"||imgur.com^",
	}},
	{"dailymotion", []string{
		"||dailymotion.com^",
		"||dm-event.net^",
		"||dmcdn.net^",
	}},
	{"qq", []string{
		// block qq.com and subdomains excluding WeChat domains
		"||qq.com^$denyallow=wx*.qq.com|weixin.qq.com",
		"||qqzaixian.com^",
	}},
	{"wechat", []string{
		"||wechat.com^",
		"||weixin.qq.com^",
		"||wx.qq.com^",
	}},
	{"viber", []string{
		"||viber.com^",
	}},
	{"weibo", []string{
		"||weibo.com^",
	}},
	{"9gag", []string{
		"||9cache.com^",
		"||9gag.com^",
	}},
	{"telegram", []string{
		"||t.me^",
		"||telegram.me^",
		"||telegram.org^",
	}},
	{"disneyplus", []string{
		"||disney-plus.net^",
		"||disneyplus.com^",
		"||disney.playback.edge.bamgrid.com^",
		"||media.dssott.com^",
	}},
	{"hulu", []string{
		"||hulu.com^",
	}},
	{"spotify", []string{
		"/_spotify-connect._tcp.local/",
		"||spotify.com^",
		"||scdn.co^",
		"||spotify.com.edgesuite.net^",
		"||spotify.map.fastly.net^",
		"||spotify.map.fastlylb.net^",
		"||spotifycdn.net^",
		"||audio-ak-spotify-com.akamaized.net^",
		"||audio4-ak-spotify-com.akamaized.net^",
		"||heads-ak-spotify-com.akamaized.net^",
		"||heads4-ak-spotify-com.akamaized.net^",
	}},
	{"tinder", []string{
		"||gotinder.com^",
		"||tinder.com^",
		"||tindersparks.com^",
	}},
}

// convert array to map
func initBlockedServices() {
	serviceRules = make(map[string][]*rules.NetworkRule)
	for _, s := range serviceRulesArray {
		netRules := []*rules.NetworkRule{}
		for _, text := range s.rules {
			rule, err := rules.NewNetworkRule(text, BlockedSvcsListID)
			if err != nil {
				log.Error("rules.NewNetworkRule: %s  rule: %s", err, text)
				continue
			}
			netRules = append(netRules, rule)
		}
		serviceRules[s.name] = netRules
	}
}

// BlockedSvcKnown - return TRUE if a blocked service name is known
func BlockedSvcKnown(s string) bool {
	_, ok := serviceRules[s]
	return ok
}

// ApplyBlockedServices - set blocked services settings for this DNS request
func (d *DNSFilter) ApplyBlockedServices(setts *Settings, list []string, global bool) {
	setts.ServicesRules = []ServiceEntry{}
	if global {
		d.confLock.RLock()
		defer d.confLock.RUnlock()
		list = d.Config.BlockedServices
	}
	for _, name := range list {
		rules, ok := serviceRules[name]

		if !ok {
			log.Error("unknown service name: %s", name)
			continue
		}

		s := ServiceEntry{}
		s.Name = name
		s.Rules = rules
		setts.ServicesRules = append(setts.ServicesRules, s)
	}
}

func (d *DNSFilter) handleBlockedServicesList(w http.ResponseWriter, r *http.Request) {
	d.confLock.RLock()
	list := d.Config.BlockedServices
	d.confLock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(list)
	if err != nil {
		aghhttp.Error(r, w, http.StatusInternalServerError, "json.Encode: %s", err)

		return
	}
}

func (d *DNSFilter) handleBlockedServicesSet(w http.ResponseWriter, r *http.Request) {
	list := []string{}
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		aghhttp.Error(r, w, http.StatusBadRequest, "json.Decode: %s", err)

		return
	}

	d.confLock.Lock()
	d.Config.BlockedServices = list
	d.confLock.Unlock()

	log.Debug("Updated blocked services list: %d", len(list))

	d.ConfigModified()
}

// registerBlockedServicesHandlers - register HTTP handlers
func (d *DNSFilter) registerBlockedServicesHandlers() {
	d.Config.HTTPRegister(http.MethodGet, "/control/blocked_services/list", d.handleBlockedServicesList)
	d.Config.HTTPRegister(http.MethodPost, "/control/blocked_services/set", d.handleBlockedServicesSet)
}
