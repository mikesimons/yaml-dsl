MRuby::Build.new do |conf|
  # Gets set by the VS command prompts.
  if ENV['VisualStudioVersion'] || ENV['VSINSTALLDIR']
    toolchain :visualcpp
  else
    toolchain :gcc
  end

  enable_debug

  conf.gembox 'default'

  # See https://github.com/mruby/mruby/blob/master/doc/guides/mrbgems.md for more about mrbgems

  #conf.gem :git => 'https://github.com/iij/mruby-regexp-pcre.git'
end