Backbone.emulateHTTP = true;
Backbone.emulateJSON = true;

var NEWS_ITEM = Backbone.Model.extend({
  defaults : {
    id : 0,
    title : "",
    img : ""
  },
  url : function(){
    return "/api/news/" + this.get('id')
  },

  get_img : function(){

    var reg = new RegExp(/style=\"(.*)\"/g)

    if( !reg.test(this.get('img')) ){
      return $('<div/>').append(
        $('<img/>').attr('src',this.get('img'))
      ).html()
    }
  }
})

var NEWS = Backbone.Collection.extend({
  model : NEWS_ITEM,
  url : "/api/news/"
})


var RULE = Backbone.Model.extend({
  defaults : {
    id : 0,
    name : "",
    link : "",
    main_path : "",
    img_path : "",
    img_attr : "",
    title_path : "",
    href_path : "",
    desc_path : "",
  },

  url : function(){
    return "/api/rule/" + this.get('id')
  },


  validate : function(){
    if( this.get('name') == this.defaults.name ){
      return "name empty"
    }
    if( this.get('link') == this.defaults.link ){
      return "link empty"
    }
    if( this.get('main_path') == this.defaults.main_path ){
      return "main_path empty"
    }
    if( this.get('img_path') == this.defaults.img_path ){
      return "img_path empty"
    }
    if( this.get('title_path') == this.defaults.title_path ){
      return "title_path empty"
    }
    if( this.get('desc_path') == this.defaults.desc_path ){
      return "desc_path empty"
    }
  }
})


var RULES = Backbone.Collection.extend({
  model : RULE,
  url : "/api/rule/"
})

var MODAL_MORE_NEWS = Backbone.View.extend({
  template : "#more_view_news",

  initialize : function( news ){
    this.model = news
    this.render()
  },

  render : function(){
    var tpl_c = _.template( $(this.template).html() )
    this.$el = $(
      tpl_c({
        news : this.model
      })
    ).modal('show')
  }
})


var NEWS_LIST = Backbone.View.extend({
  el : "#news_list",

  template : "#news_list_tpl",

  collection : new NEWS,

  initialize: function( news ){
    this.collection = news
    this.listenTo( this.collection , "sync" , this.render )
    this.render()
  },

  render : function(){
    var tpl_c = _.template( $(this.template).html())
    $(this.$el).empty().append(
      tpl_c({
        news : this.collection
      })
    )
  },

  events : {
    "click .more_view_news" : function( e ){
      var cid = $(e.target).attr('id')
      var model = this.collection.get(cid)
      if( typeof model != "undefined" ){
        new MODAL_MORE_NEWS( model )
      }
    }
  }
})


var RULE_MODAL = Backbone.View.extend({
  template : "#rule_modal_tpl",

  initialize : function( rule , rules ){

    this.model = rule
    this.collection = rules

    this.listenTo( this.model , "change" , this.__btn_save_e)

    this.render()

  },

  __btn_save_e : function(){
    if( this.model.isValid() ){
      $("#save").removeClass('btn-disabled').addClass('btn-success')
    }else{
      $("#save").addClass('btn-disabled').removeClass('btn-success')
    }
  },

  render : function(){
    var self = this
    var tpl_c = _.template( $(this.template).html() )
    this.$el = $(
      tpl_c({
        rule : this.model
      })
    ).modal('show')

    $(this.$el).on('shown.bs.modal' , function(){
      self.__btn_save_e()
    }).on('hidden.bs.modal' , function(){
      $(this).remove()
      self.collection.fetch()
    })
  },

  events : {
    "keyup #name" : function( e ){
      this.model.set('name' , $(e.target).val() )
    },
    "keyup #link" : function( e ){
      this.model.set('link' , $(e.target).val() )
    },
    "keyup #main_path" : function( e ){
      this.model.set('main_path' , $(e.target).val() )
    },
    "keyup #img_path" : function( e ){
      this.model.set('img_path' , $(e.target).val() )
    },
    "keyup #img_attr" : function( e ){
      this.model.set('img_attr' , $(e.target).val() )
    },
    "keyup #title_path" : function( e ){
      this.model.set('title_path' , $(e.target).val() )
    },
    "keyup #href_path" : function( e ){
      this.model.set('href_path' , $(e.target).val() )
    },
    "keyup #desc_path" : function( e ){
      this.model.set('desc_path' , $(e.target).val() )
    },
    "click #is_img_attr" : function( e ){
      var action = $(e.target).prop('checked') ? "show" : "hide"
      $('#img_attr').closest('.form-group')[action]()
    },
    "click #is_blank" : function( e ){
      var action = $(e.target).prop('checked') ? "show" : "hide"
      $('#href_path').closest('.form-group')[action]()
    },
    "click #save" : function( e ){
      this.model.save()
      $(this.$el).modal('hide')
    }
  }
})

var RULE_MODAL_DEL = Backbone.View.extend({

  template : "#rules_delete_modal",

  initialize : function( model , collection ){

    this.model = model
    this.collection = collection

    this.render()

  },
  render : function(){
    var self = this
    var tpl_c = _.template( $(this.template).html() )
    this.$el = $(
      tpl_c({
        rule : this.model
      })
    ).modal('show')

    $(this.$el).on('hidden.bs.modal' , function(){
      $(this).remove()
      self.collection.fetch()
    })
  },

  events : {
    "click #delete" : function(){
      this.model.destroy()
      $(this.$el).modal('hide')
    }
  }
})



var RULES_LIST = Backbone.View.extend({
  el : "#list_rules",

  template : "#list_rules_tpl",

  collection : new RULES,

  initialize: function( news ){
    this.collection = news
    this.listenTo( this.collection , "sync" , this.render )
    this.render()
  },

  render : function(){
    var tpl_c = _.template( $(this.template).html())
    $(this.$el).empty().append(
      tpl_c({
        news : this.collection
      })
    )
  },

  events : {
    "click #add_rule" : function(){
      new RULE_MODAL( new RULE  , this.collection)
    },

    "click .delete_rule" : function(e){
      var cid = $(e.target).closest('li').attr('id')
      var model = this.collection.get(cid)
      if( typeof model != "undefined" ){
        new RULE_MODAL_DEL( model  , this.collection)
      }
    },

    "click .edit_rule" : function(e){
      var cid = $(e.target).closest('li').attr('id')
      var model = this.collection.get(cid)
      if( typeof model != "undefined" ){
        new RULE_MODAL( model  , this.collection)
      }
    }
  }
})


var news_collection = new NEWS()

var rules = new RULES
rules.fetch()
news_collection.fetch ()

$(document).ready(function(){
  new NEWS_LIST( news_collection )
  new RULES_LIST( rules )
})