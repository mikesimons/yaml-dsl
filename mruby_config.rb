MRuby::Build.new do |conf|
  # Gets set by the VS command prompts.
  if ENV['VisualStudioVersion'] || ENV['VSINSTALLDIR']
    toolchain :visualcpp
  else
    toolchain :gcc
  end

  enable_debug

  conf.gembox 'full-core'
  conf.gem :git => 'https://github.com/jbreeden/mruby-erb.git'
  conf.gem :git => 'https://github.com/ksss/mruby-ostruct.git'
  conf.gem :git => 'https://github.com/AndrewBelt/mruby-yaml.git'
  
  # See https://github.com/mruby/mruby/blob/master/doc/guides/mrbgems.md for more about mrbgems

  #conf.gem :git => 'https://github.com/iij/mruby-regexp-pcre.git'
end
