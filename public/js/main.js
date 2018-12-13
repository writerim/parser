var partition_size = 15

//////////////////////
//   MODELS
//////////////////////


// Этот вид отвечает за готовность всех моделей
var READY_STAT = Backbone.Model.extend({
  defaults : {
    news_collect : false,
    sources_collect : false,
    news_attrs_value_collect : false,
    news_attrs_collection : false,
    rules : false
  },
  is_ready : function(){
    return this.get('news_collect') &&
    this.get('sources_collect') &&
    this.get('news_attrs_value_collect') && 
    this.get('news_attrs_collection') &&
    this.get('rules')
  }
})


var RULE = Backbone.Model.extend({
  defaults : {
    id : 0,
    news_attr_id : 0,
    source_list_id : 0,
    rule : "",
    get_attr : "",
    is_main : 0,
    is_unique : 0
  }
})

var NEWS = Backbone.Model.extend({
  defaults : {
    id : 0,
    source_id : 0,
    title : ""
  }
})

var NEWS_ATTRS_VALUE = Backbone.Model.extend({
  defaults : {
    id : 0,
    news_id : 0,
    news_attrs_id : 0,
    value : ""
  }
})

var NEWS_ATTRS = Backbone.Model.extend({
  defaults : {
    id : 0,
    ident : "",
    name : ""
  }
})

var SOURCE_MODEL = Backbone.Model.extend({
  defaults : {
    id : 0,
    name : "",
    href : ""
  }
})



//////////////////////
//   COLLECTIONS
//////////////////////

var SOURCE_COLLECTION = Backbone.Collection.extend({
  model : SOURCE_MODEL,
  url : "/api/source/0/1000/" // TODO Пока так
})

var RULES = Backbone.Collection.extend({
  model : RULE,
  url : "/api/news_rule_list/"
})

var NEWS_ATTRS_COLLECTION = Backbone.Collection.extend({
  model : NEWS_ATTRS,
  url : "/api/news_attrs/"
})

var NEWS_ATTRS_VALUE_COLLECTION = Backbone.Collection.extend({
  model : NEWS_ATTRS_VALUE,
  url : function(){
    return "/api/news_attrs_value/0/" + ( this.length + partition_size ) + "/"
  }
})

var LIST_NEWS = Backbone.Collection.extend({ 
  model : NEWS,
  url : function(){
    return "/api/news/0/" + ( news_collect.length + partition_size ) + "/"
  }
})








//////////////////////
//   VIEWS
//////////////////////



var LIST_RULES = Backbone.View.extend({
  el : "#rules",

  template : "#rules_tpl",

  collection : new SOURCE_COLLECTION,

  initialize : function( collection ){
    this.collection = collection
    this.listenTo( this.collection , "sync" , this.render )
    this.render()
  },

  render : function(){
    var tpl_c = _.template($(this.template).html())
    $(this.$el).empty().append(tpl_c({
      sources : this.collection
    }))
  }
})



var LIST_NEWS_VIEW = Backbone.View.extend({
  el : "#list_news",

  template : "#list_news_tpl",

  collection : new LIST_NEWS,

  initialize : function( news , attrs , attrs_value, rules , sources ){
    this.collection = news
    this.attrs = attrs
    this.attrs_value = attrs_value
    this.rules = rules
    this.sources = sources
    this.listenTo( this.collection , "sync" , this.render )
    this.listenTo( this.attrs , "sync" , this.render )
    this.listenTo( this.attrs_value , "sync" , this.render )
    this.listenTo( this.rules , "sync" , this.render )
    this.listenTo( this.sources , "sync" , this.render )
    this.render()
  },

  render : function(){
    var tpl_c = _.template($(this.template).html())
    $(this.$el).empty().append(tpl_c({
      news : this.collection,
      attrs : this.attrs,
      rules : this.rules,
      sources : this.sources,
      attrs_value : this.attrs_value
    }))
  },

  events : {
    "click #more_news" : function(){
      var self = this
      this.collection.fetch({ 
        success : function(){
          self.attrs_value.fetch({ url : "/api/news_attrs_value/0/" + ( self.collection.length + partition_size ) + "/" })
        } 
      })
    }
  }
})



//////////////////////
//   BODY
//////////////////////



var news_collect = new LIST_NEWS
var sources_collect = new SOURCE_COLLECTION
var news_attrs_value_collect = new NEWS_ATTRS_VALUE_COLLECTION
var news_attrs_collection = new NEWS_ATTRS_COLLECTION
var rules = new RULES


var ready_stat = new READY_STAT

$(document).ready(function(){

  sources_collect.fetch({ success : function(){ ready_stat.set('sources_collect', true ) }})
  news_collect.fetch({ success : function(){  ready_stat.set('news_collect', true ) } })
  news_attrs_collection.fetch({ success : function(){  ready_stat.set('news_attrs_collection', true )  } })
  news_attrs_value_collect.fetch({ success : function(){  ready_stat.set('news_attrs_value_collect', true )  } })
  rules.fetch({ success : function(){  ready_stat.set('rules', true )  } })

  ready_stat.on('change' , function(){
    if( this.is_ready() ){
      new LIST_NEWS_VIEW( news_collect , news_attrs_collection ,news_attrs_value_collect , rules , sources_collect )
      new LIST_RULES( sources_collect )
    }
  })

})