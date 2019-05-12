module Jekyll
  class EpisodePageGenerator < Generator
    def generate(site)
      site.data['index'].each do |feed, episodes|
        episodes.each do |episode|
          site.pages << EpisodePage.new(site, site.source, "#{feed}", episode)
        end
      end
    end
  end

  class EpisodePage < Page
    def initialize(site, base, dir, episode)
      @site = site
      @base = base
      @dir = dir
      @name = "#{episode['number']}-#{Jekyll::Utils.slugify(episode['title'])}.html"

      process(name)
      read_yaml('_layouts', 'episode.html')
      self.data['transcript'] = File.read("_data/episodes/#{episode['source']}").strip.split("\n")
    end
  end
end
