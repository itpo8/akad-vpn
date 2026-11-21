#!/usr/bin/env ruby
# Removes the Packet Tunnel *only from a temporary CI checkout* to produce a
# UI-only iOS demo that can be signed with an ordinary Apple ID.
# This is not a VPN-capable build. Never use this as a VPN release.
require 'xcodeproj'

path = File.expand_path('../app/ios/Runner.xcodeproj', __dir__)
project = Xcodeproj::Project.open(path)
runner = project.targets.find { |target| target.name == 'Runner' }
abort('Runner target missing') unless runner

extension = project.targets.find { |target| target.name == 'AkadPacketTunnel' }
if extension
  runner.dependencies.select { |item| item.target == extension }.each(&:remove_from_project)
end

# Remove any extension embeds so Xcode does not demand Network Extension
# provisioning for this visual-only preview.
runner.build_phases.select do |phase|
  phase.is_a?(Xcodeproj::Project::Object::PBXCopyFilesBuildPhase) &&
    phase.name == 'Embed App Extensions'
end.each(&:remove_from_project)

runner.build_configurations.each do |configuration|
  configuration.build_settings.delete('CODE_SIGN_ENTITLEMENTS')
end

project.save
puts 'AkadVPN visual-only iOS demo: VPN extension removed from Runner build.'
