package config

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/m1k1o/neko/server/internal/connectivity"
	"github.com/m1k1o/neko/server/pkg/types"
	"github.com/m1k1o/neko/server/pkg/utils"
)

// default stun server
const defStunSrv = "stun:stun.l.google.com:19302"

type WebRTCEstimator struct {
	Enabled        bool
	Passive        bool
	Debug          bool
	InitialBitrate int

	// how often to read and process bandwidth estimation reports
	ReadInterval time.Duration
	// how long to wait for stable connection (only neutral or upward trend) before upgrading
	StableDuration time.Duration
	// how long to wait for unstable connection (downward trend) before downgrading
	UnstableDuration time.Duration
	// how long to wait for stalled connection (neutral trend with low bandwidth) before downgrading
	StalledDuration time.Duration
	// how long to wait before downgrading again after previous downgrade
	DowngradeBackoff time.Duration
	// how long to wait before upgrading again after previous upgrade
	UpgradeBackoff time.Duration
	// how bigger the difference between estimated and stream bitrate must be to trigger upgrade/downgrade
	DiffThreshold float64
}

type WebRTC struct {
	ICELite            bool
	ICETrickle         bool
	ICEServersFrontend []types.ICEServer
	ICEServersBackend  []types.ICEServer
	TCPMux             int
	UDPMux             int

	NAT1To1IPs     []string
	IpRetrievalUrl string
	Connectivity   connectivity.Mode

	Estimator WebRTCEstimator
}

func (WebRTC) Init(cmd *cobra.Command) error {
	cmd.PersistentFlags().Bool("webrtc.icelite", false, "configures whether or not the ICE agent should be a lite agent")
	if err := viper.BindPFlag("webrtc.icelite", cmd.PersistentFlags().Lookup("webrtc.icelite")); err != nil {
		return err
	}

	cmd.PersistentFlags().Bool("webrtc.icetrickle", true, "configures whether cadidates should be sent asynchronously using Trickle ICE")
	if err := viper.BindPFlag("webrtc.icetrickle", cmd.PersistentFlags().Lookup("webrtc.icetrickle")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("webrtc.iceservers.frontend", "[]", "STUN and TURN servers used by the frontend")
	if err := viper.BindPFlag("webrtc.iceservers.frontend", cmd.PersistentFlags().Lookup("webrtc.iceservers.frontend")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("webrtc.iceservers.backend", "[]", "STUN and TURN servers used by the backend")
	if err := viper.BindPFlag("webrtc.iceservers.backend", cmd.PersistentFlags().Lookup("webrtc.iceservers.backend")); err != nil {
		return err
	}

	cmd.PersistentFlags().Int("webrtc.tcpmux", 0, "single TCP mux port for all peers")
	if err := viper.BindPFlag("webrtc.tcpmux", cmd.PersistentFlags().Lookup("webrtc.tcpmux")); err != nil {
		return err
	}

	cmd.PersistentFlags().Int("webrtc.udpmux", 52000, "single UDP mux port for all peers")
	if err := viper.BindPFlag("webrtc.udpmux", cmd.PersistentFlags().Lookup("webrtc.udpmux")); err != nil {
		return err
	}

	cmd.PersistentFlags().StringSlice("webrtc.nat1to1", []string{}, "sets a list of external IP addresses of 1:1 (D)NAT and a candidate type for which the external IP address is used")
	if err := viper.BindPFlag("webrtc.nat1to1", cmd.PersistentFlags().Lookup("webrtc.nat1to1")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("webrtc.ip_retrieval_url", "https://checkip.amazonaws.com", "URL address used for retrieval of the external IP address")
	if err := viper.BindPFlag("webrtc.ip_retrieval_url", cmd.PersistentFlags().Lookup("webrtc.ip_retrieval_url")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("webrtc.connectivity.mode", "direct", "connectivity mode (direct, frp, or auto)")
	if err := viper.BindPFlag("webrtc.connectivity.mode", cmd.PersistentFlags().Lookup("webrtc.connectivity.mode")); err != nil {
		return err
	}

	// bandwidth estimator

	cmd.PersistentFlags().Bool("webrtc.estimator.enabled", false, "enables the bandwidth estimator")
	if err := viper.BindPFlag("webrtc.estimator.enabled", cmd.PersistentFlags().Lookup("webrtc.estimator.enabled")); err != nil {
		return err
	}

	cmd.PersistentFlags().Bool("webrtc.estimator.passive", false, "passive estimator mode, when it does not switch pipelines, only estimates")
	if err := viper.BindPFlag("webrtc.estimator.passive", cmd.PersistentFlags().Lookup("webrtc.estimator.passive")); err != nil {
		return err
	}

	cmd.PersistentFlags().Bool("webrtc.estimator.debug", false, "enables debug logging for the bandwidth estimator")
	if err := viper.BindPFlag("webrtc.estimator.debug", cmd.PersistentFlags().Lookup("webrtc.estimator.debug")); err != nil {
		return err
	}

	cmd.PersistentFlags().Int("webrtc.estimator.initial_bitrate", 1_000_000, "initial bitrate for the bandwidth estimator")
	if err := viper.BindPFlag("webrtc.estimator.initial_bitrate", cmd.PersistentFlags().Lookup("webrtc.estimator.initial_bitrate")); err != nil {
		return err
	}

	cmd.PersistentFlags().Duration("webrtc.estimator.read_interval", 2*time.Second, "how often to read and process bandwidth estimation reports")
	if err := viper.BindPFlag("webrtc.estimator.read_interval", cmd.PersistentFlags().Lookup("webrtc.estimator.read_interval")); err != nil {
		return err
	}

	cmd.PersistentFlags().Duration("webrtc.estimator.stable_duration", 12*time.Second, "how long to wait for stable connection (upward or neutral trend) before upgrading")
	if err := viper.BindPFlag("webrtc.estimator.stable_duration", cmd.PersistentFlags().Lookup("webrtc.estimator.stable_duration")); err != nil {
		return err
	}

	cmd.PersistentFlags().Duration("webrtc.estimator.unstable_duration", 6*time.Second, "how long to wait for stalled connection (neutral trend with low bandwidth) before downgrading")
	if err := viper.BindPFlag("webrtc.estimator.unstable_duration", cmd.PersistentFlags().Lookup("webrtc.estimator.unstable_duration")); err != nil {
		return err
	}

	cmd.PersistentFlags().Duration("webrtc.estimator.stalled_duration", 24*time.Second, "how long to wait for stalled bandwidth estimation before downgrading")
	if err := viper.BindPFlag("webrtc.estimator.stalled_duration", cmd.PersistentFlags().Lookup("webrtc.estimator.stalled_duration")); err != nil {
		return err
	}

	cmd.PersistentFlags().Duration("webrtc.estimator.downgrade_backoff", 10*time.Second, "how long to wait before downgrading again after previous downgrade")
	if err := viper.BindPFlag("webrtc.estimator.downgrade_backoff", cmd.PersistentFlags().Lookup("webrtc.estimator.downgrade_backoff")); err != nil {
		return err
	}

	cmd.PersistentFlags().Duration("webrtc.estimator.upgrade_backoff", 5*time.Second, "how long to wait before upgrading again after previous upgrade")
	if err := viper.BindPFlag("webrtc.estimator.upgrade_backoff", cmd.PersistentFlags().Lookup("webrtc.estimator.upgrade_backoff")); err != nil {
		return err
	}

	cmd.PersistentFlags().Float64("webrtc.estimator.diff_threshold", 0.15, "how bigger the difference between estimated and stream bitrate must be to trigger upgrade/downgrade")
	if err := viper.BindPFlag("webrtc.estimator.diff_threshold", cmd.PersistentFlags().Lookup("webrtc.estimator.diff_threshold")); err != nil {
		return err
	}

	return nil
}

func (s *WebRTC) Set() {
	if err := validateUnsupportedLegacyConfig(); err != nil {
		log.Panic().Err(err).Msg("unsupported WebRTC configuration")
	}

	s.ICELite = viper.GetBool("webrtc.icelite")
	s.ICETrickle = viper.GetBool("webrtc.icetrickle")

	// parse frontend ice servers
	if err := viper.UnmarshalKey("webrtc.iceservers.frontend", &s.ICEServersFrontend, viper.DecodeHook(
		utils.JsonStringAutoDecode(s.ICEServersFrontend),
	)); err != nil {
		log.Warn().Err(err).Msgf("unable to parse frontend ICE servers")
	}

	// parse backend ice servers
	if err := viper.UnmarshalKey("webrtc.iceservers.backend", &s.ICEServersBackend, viper.DecodeHook(
		utils.JsonStringAutoDecode(s.ICEServersBackend),
	)); err != nil {
		log.Warn().Err(err).Msgf("unable to parse backend ICE servers")
	}

	if s.ICELite && len(s.ICEServersBackend) > 0 {
		log.Warn().Msgf("ICE Lite is enabled, but backend ICE servers are configured. Backend ICE servers will be ignored.")
	}

	// Use the documented default STUN server only when neither endpoint has an
	// explicit ICE server. The obsolete global webrtc.iceservers value is
	// deliberately rejected before this point.
	if len(s.ICEServersFrontend) == 0 && len(s.ICEServersBackend) == 0 {
		defaultServer := types.ICEServer{URLs: []string{defStunSrv}}
		s.ICEServersFrontend = append(s.ICEServersFrontend, defaultServer)
		s.ICEServersBackend = append(s.ICEServersBackend, defaultServer)
	}

	s.TCPMux = viper.GetInt("webrtc.tcpmux")
	s.UDPMux = viper.GetInt("webrtc.udpmux")

	s.Connectivity = connectivity.Mode(viper.GetString("webrtc.connectivity.mode"))
	s.NAT1To1IPs = viper.GetStringSlice("webrtc.nat1to1")
	s.IpRetrievalUrl = viper.GetString("webrtc.ip_retrieval_url")
	if s.Connectivity == connectivity.ModeAuto {
		// A server behind NAT must gather its mapped address from STUN. The
		// HTTP endpoint only reports the NAT gateway/CGNAT address and would
		// incorrectly advertise it as a static 1:1 host candidate.
		s.IpRetrievalUrl = ""
	}
	if s.IpRetrievalUrl != "" && len(s.NAT1To1IPs) == 0 {
		ip, err := utils.HttpRequestGET(s.IpRetrievalUrl)
		if err == nil {
			s.NAT1To1IPs = append(s.NAT1To1IPs, ip)
		} else {
			log.Warn().Err(err).Msgf("IP retrieval failed")
		}
	}

	if err := s.validateConnectivity(); err != nil {
		log.Panic().Err(err).Msg("invalid WebRTC connectivity configuration")
	}

	// bandwidth estimator

	s.Estimator.Enabled = viper.GetBool("webrtc.estimator.enabled")
	s.Estimator.Passive = viper.GetBool("webrtc.estimator.passive")
	s.Estimator.Debug = viper.GetBool("webrtc.estimator.debug")
	s.Estimator.InitialBitrate = viper.GetInt("webrtc.estimator.initial_bitrate")
	s.Estimator.ReadInterval = viper.GetDuration("webrtc.estimator.read_interval")
	s.Estimator.StableDuration = viper.GetDuration("webrtc.estimator.stable_duration")
	s.Estimator.UnstableDuration = viper.GetDuration("webrtc.estimator.unstable_duration")
	s.Estimator.StalledDuration = viper.GetDuration("webrtc.estimator.stalled_duration")
	s.Estimator.DowngradeBackoff = viper.GetDuration("webrtc.estimator.downgrade_backoff")
	s.Estimator.UpgradeBackoff = viper.GetDuration("webrtc.estimator.upgrade_backoff")
	s.Estimator.DiffThreshold = viper.GetFloat64("webrtc.estimator.diff_threshold")
}

func (s WebRTC) validateConnectivity() error {
	if s.Connectivity == connectivity.ModeFRP {
		if !viper.IsSet("webrtc.nat1to1") {
			return fmt.Errorf("frp connectivity mode requires an explicit webrtc.nat1to1 IP")
		}
		if len(s.NAT1To1IPs) != 1 {
			return fmt.Errorf("frp connectivity mode requires exactly one webrtc.nat1to1 IP")
		}

		return connectivity.MediaPortPlan{
			Mode:       s.Connectivity,
			UDPMuxPort: s.UDPMux,
			TCPMuxPort: s.TCPMux,
			NAT1To1IP:  s.NAT1To1IPs[0],
		}.Validate()
	}

	if s.Connectivity == connectivity.ModeDirect {
		plan := connectivity.MediaPortPlan{
			Mode:       s.Connectivity,
			UDPMuxPort: s.UDPMux,
			TCPMuxPort: s.TCPMux,
		}
		if len(s.NAT1To1IPs) == 1 {
			plan.NAT1To1IP = s.NAT1To1IPs[0]
		}
		return plan.Validate()
	}

	if s.Connectivity == connectivity.ModeAuto {
		if s.ICELite {
			return fmt.Errorf("auto connectivity mode requires webrtc.icelite=false")
		}
		if len(s.NAT1To1IPs) != 0 {
			return fmt.Errorf("auto connectivity mode must not set webrtc.nat1to1")
		}
		if len(s.ICEServersFrontend) == 0 || len(s.ICEServersBackend) == 0 {
			return fmt.Errorf("auto connectivity mode requires frontend and backend ICE servers")
		}
		return connectivity.MediaPortPlan{Mode: s.Connectivity}.Validate()
	}

	return fmt.Errorf("unsupported connectivity mode %q", s.Connectivity)
}

func validateUnsupportedLegacyConfig() error {
	if viper.IsSet("webrtc.epr") {
		return fmt.Errorf("webrtc.epr is no longer supported; configure webrtc.udpmux (and optionally webrtc.tcpmux) instead")
	}

	if usesLegacyGlobalICEServers(viper.Get("webrtc.iceservers")) {
		return fmt.Errorf("webrtc.iceservers is no longer supported; configure webrtc.iceservers.frontend and webrtc.iceservers.backend instead")
	}

	return nil
}

func usesLegacyGlobalICEServers(value any) bool {
	if value == nil {
		return false
	}

	// A nested frontend/backend configuration is represented as a map. Any
	// other value is the old shared ICE-server list (including its JSON
	// environment-variable form).
	if values, ok := value.(map[string]any); ok {
		for key := range values {
			if key != "frontend" && key != "backend" {
				return true
			}
		}
		return false
	}

	return true
}
