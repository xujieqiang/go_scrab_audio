# go_scrab_audio
go  colly 爬虫     获取音频


采用colly
爬虫练习。
音频爬取的时候，是按照音频所在资源分布规律分析出来的。
每次用  
c.OnHTML("div[id=fdfarrer]"，func(e *colly.HTMLElement){
  audio:=e.Attr("src")
  
})


并不能获取 audio 项目的src  链接， 
好像这是动态数据，

结果：：能把所有的音频全部爬下来。



**仍然有些问题 没有解决，不能爬取动态数据
