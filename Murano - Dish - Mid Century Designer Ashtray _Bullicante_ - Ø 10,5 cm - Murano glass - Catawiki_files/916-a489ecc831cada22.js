(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[916],{80607:function(e,n,t){"use strict";t.d(n,{n:function(){return a}});var r=t(27267),a=Object.assign((function(e){var n=e.refs,t=e.onClose;return function(e){var a=n.floating.current,u=n.reference.current;null===a||e.currentTarget===a||e.currentTarget===u||function(e,n){var t=n.getBoundingClientRect(),a=(0,r.Z)(e.clientX,t.left,t.right),u=(0,r.Z)(e.clientY,t.top-10,t.bottom);return a&&u}(e,a)||t()}}),{__options:{blockPointerEvents:!0}})},21389:function(e,n,t){"use strict";t.d(n,{$y:function(){return u},G1:function(){return s},dR:function(){return i},f9:function(){return c},l7:function(){return l},sR:function(){return o}});var r=t(68697),a=function(e){return e.lotBidSubmitStatus},u=(0,r.P1)(a,(function(e){return e.bidResponseType})),i=(0,r.P1)(a,(function(e){return e.description})),c=(0,r.P1)(a,(function(e){return e.title})),o=(0,r.P1)(a,(function(e){return e.adviceType})),s=(0,r.P1)(a,(function(e){return e.isBidConfirmed})),l=(0,r.P1)(a,(function(e){return e.isBidSubmitStatusDialogueVisible}))},45247:function(e,n,t){"use strict";t.d(n,{K:function(){return k}});var r=t(85893),a=t(67294),u=t(2664),i=t(5309),c=t(93967),o=t.n(c),s=t(73391),l=t(83493),d=t(91448),v=t(3173),f=t(58072),h=t(63286),p=t(34752),g=t(7755),_=t(47240),y=t(47568),b=t(97582),m=t(43240),S=t(3735),T=t(30365),C=t(71835),L=t(23521),A=t(36551),x=t(7842),V=function(e){var n=e.isOpen,t=e.searchTerm,u=e.anchor,i=e.handleAddTerm,c=(0,s.QT)(),o=(0,a.useState)([]),l=o[0],d=o[1];return(0,a.useEffect)((function(){var e=function(){var e=(0,y.Z)((function(){return(0,b.__generator)(this,(function(e){switch(e.label){case 0:return[4,(0,x.KF)(t,c.locale)];case 1:return S.zG.apply(void 0,[e.sent(),(0,T.match)((function(e){(0,A.WB)("savedSearchesSuggestions")(e)}),(function(e){d(e)}))]),[2]}}))}));return function(){return e.apply(this,arguments)}}();(null===t||void 0===t?void 0:t.length)>=_.GA&&e()}),[t,c.locale]),(0,r.jsx)(C.v,{open:n,onToggle:m.Z,anchor:u,focusOrder:["reference"],children:null===l||void 0===l?void 0:l.map((function(e){return(0,r.jsx)(L.c,{onClick:i(e.text),children:e.text},e.text)}))})},j=t(94672),B=t.n(j),k=function(e){var n=e.inputValue,t=e.suggestionAdded,c=e.setSuggestionAdded,y=e.setInputValue,b=e.handleAddTerm,m=e.className,S=e.textFieldClassName,T=e.discardEdittingSearchTerm,C=e.withSubmitButton,L=void 0!==C&&C,A=e.submitButtonText,x=e.submitButtonWrapperClassName,j=(0,s.QT)(),k=(0,h.C)(p.eD),E=(0,u.v9)((function(e){return(0,f.c7)(e)})),Z=(0,a.useState)(!1),D=Z[0],q=Z[1],w=(0,a.useRef)(null);(0,i.O)(w,(function(){T()})),(0,a.useEffect)((function(){t||q(n.length>=_.GA)}),[n.length,t]),(0,a.useEffect)((function(){t&&q(!1)}),[t]),(0,a.useEffect)((function(){var e,t;n&&(null===(t=null===(e=w.current)||void 0===e?void 0:e.querySelector("input"))||void 0===t||t.focus())}),[n]);return(0,r.jsx)("div",{ref:w,onKeyDown:function(e){"Enter"===e.key&&n&&(e.preventDefault(),b(n),(0,g.Xc)(n,E))},className:o()(B().container,m),children:(0,r.jsxs)("div",{className:"u-flex",children:[(0,r.jsx)("div",{className:"u-flex--1",children:(0,r.jsx)(V,{isOpen:D,searchTerm:n,handleAddTerm:function(e){return function(){b(e),(0,g.BC)(e,E)}},anchor:(0,r.jsx)(l.n,{inputClassName:B().input,value:n,onChange:function(e){y(e.currentTarget.value),null===c||void 0===c||c(!1)},placeholder:j.t("add_new_search"),autoComplete:"off",autoCapitalize:"off",type:"search",className:S})})}),L&&!k&&(0,r.jsx)("div",{className:x,children:(0,r.jsx)(d.z,{onClick:function(){n&&(b(n),(0,g.Xc)(n,E))},type:"button",icon:v.F,iconPosition:"start",className:"u-no-wrap",children:A||j.t("interests.search.save")})})]})})}},47240:function(e,n,t){"use strict";t.d(n,{AG:function(){return s},CX:function(){return r},E$:function(){return u},GA:function(){return o},ZK:function(){return c},cK:function(){return i},jQ:function(){return a}});var r=1,a=10,u=6,i=800,c=10,o=3,s="recentSearches"},7842:function(e,n,t){"use strict";t.d(n,{iQ:function(){return m},cZ:function(){return L},pp:function(){return b},KF:function(){return T},lb:function(){return S},_6:function(){return C}});var r=t(14924),a=t(30902),u=t(30365),i=t(47821),c=t(80896),o=t(11564),s=t(81422),l=t(54335),d=t(36551),v=t(65096),f=t(32521),h=f.interface({created_at:f.string,id:f.number,locale:f.string,query:f.string,updated_at:f.string,user_id:f.number,match_count:c.M8.nullable(f.number)}),p=f.interface({search_alerts:f.array(h)}),g=f.interface({entity:f.interface({type:f.literal("query_term")}),highlighted:c.M8.nullable(f.string),text:f.string}),_=f.interface({query_terms:f.array(g)}),y=(f.array(f.string),t(47240)),b=function(e,n){var t;return(0,a.pipe)(c.eV.get(s.Rv),c.eV.setSearch((t={locale:n},(0,r.Z)(t,o.LZ,y.jQ),(0,r.Z)(t,o.qX,y.CX),(0,r.Z)(t,o.dY,o.Uj.PublishedAtDesc),(0,r.Z)(t,o.mo,e),(0,r.Z)(t,o.JO,!0),t)),c.eV.send()).then((function(e){return e.total?{total:e.total,lots:e.lots}:{total:0,lots:[]}}))},m=function(e,n){return(0,a.pipe)(c.eV.postJson(s.Do,{user_id:n,query:e,match_count:0}),c.eV.send()).then((function(e){return l.AB.conversion(v.in),e})).catch((function(e){return(0,d.WB)("addToSavedSearches")(e),!1}))},S=function(e,n,t){return n?m(e,t):(r=e,(0,a.pipe)(c.eV.del("".concat(s.Do,"/").concat(r))(),c.eV.send()).then((function(){return l.AB.conversion(v.qX),!0})).catch((function(e){return(0,d.WB)("deleteSavedSearch")(e),!1})));var r},T=function(e,n){return(0,a.pipe)(c.eV.get(s.iT),c.eV.setSearch({q:e,locale:n,size:y.ZK,filters:"query_terms"}),c.eV.send()).then(c.M8.decodeWith(_)).then((0,u.map)((function(e){return e.query_terms}))).catch(c.kF.errorToLeft)},C=function(e,n){return(0,a.pipe)(c.eV.patchJson("".concat(s.Do,"/").concat(e),{query:n}),c.eV.send()).then((function(e){return l.AB.conversion(v.JX),e})).catch((function(e){return(0,d.WB)("updateSavedSearch")(e),!1}))},L=function(e){return c.jT.retry((function(){return(0,a.pipe)(c.eV.get(s.Do),c.eV.setSearch({locale:e}),c.eV.send()).then(c.M8.decodeWith(p)).then((0,u.map)((function(e){var n=e.search_alerts;return(0,i.Z)(n,(function(e){return new Date(e.created_at)}),["desc"])}))).catch(c.kF.errorToLeft)}))}},7755:function(e,n,t){"use strict";t.d(n,{BC:function(){return o},DY:function(){return a},EI:function(){return u},FF:function(){return r},IA:function(){return s},XV:function(){return c},Xc:function(){return i},pb:function(){return l}});var r=function(e){return function(){var n,t;null===(n=globalThis.dataLayer)||void 0===n||n.push({event:"ga-event",eventCategory:"My Favourites events",eventAction:"Saved searches - object clicked",eventLabel:e.id}),null===(t=globalThis.dataLayer)||void 0===t||t.push({event:"my_favs_saved_searches_object_clicked",event_label:e.id})}},a=function(){var e,n;null===(e=globalThis.dataLayer)||void 0===e||e.push({event:"ga-event",eventCategory:"My Favourites events",eventAction:"Saved searches - view all clicked"}),null===(n=globalThis.dataLayer)||void 0===n||n.push({event:"my_favs_saved_searches_view_all_clicked"})},u=function(){var e,n;null===(e=globalThis.dataLayer)||void 0===e||e.push({event:"ga-event",eventCategory:"My Favourites events",eventAction:"Saved searches - scrolled"}),null===(n=globalThis.dataLayer)||void 0===n||n.push({event:"my_favs_saved_searches_scrolled"})},i=function(e,n){var t,r;null===(t=globalThis.dataLayer)||void 0===t||t.push({event:"ga-event",eventCategory:"Search events",eventAction:"Saved searches - add",eventLabel:e,user_id:n}),null===(r=globalThis.dataLayer)||void 0===r||r.push({event:"search_add_saved",event_label:e})},c=function(e,n){var t,r;null===(t=globalThis.dataLayer)||void 0===t||t.push({event:"ga-event",eventCategory:"Search events",eventAction:"Saved searches - remove",eventLabel:e,user_id:n}),null===(r=globalThis.dataLayer)||void 0===r||r.push({event:"search_remove_saved",event_label:{searchTerm:e}})},o=function(e,n){var t,r;null===(t=globalThis.dataLayer)||void 0===t||t.push({event:"ga-event",eventCategory:"Search events",eventAction:"Saved searches - add from suggestions",eventLabel:e,user_id:n}),null===(r=globalThis.dataLayer)||void 0===r||r.push({event:"search_add_saved",event_label:e,source:"Search term suggestions"})},s=function(e,n){var t,r;null===(t=globalThis.dataLayer)||void 0===t||t.push({event:"ga-event",eventCategory:"Search events",eventAction:"Saved searches - add from recent searches",eventLabel:e,user_id:n}),null===(r=globalThis.dataLayer)||void 0===r||r.push({event:"search_add_saved",event_label:e,source:"Recent searches"})},l=function(e,n){var t,r;null===(t=globalThis.dataLayer)||void 0===t||t.push({event:"ga-event",eventCategory:"Search events",eventAction:"Saved searches - remove from recent searches",eventLabel:e,user_id:n}),null===(r=globalThis.dataLayer)||void 0===r||r.push({event:"search_add_removed",event_label:e,source:"Recent searches"})}},2025:function(e,n,t){"use strict";t.d(n,{M:function(){return a}});var r=t(85893),a=(t(67294),function(e){var n=e.condition,t=e.wrapper,a=e.children;return n?t(a):(0,r.jsx)(r.Fragment,{children:a})})},34258:function(e,n,t){"use strict";t.d(n,{$$:function(){return a},Ps:function(){return u},an:function(){return r}});var r=function(){var e=globalThis.document.querySelector("header.c-header");return e?e.firstElementChild:null},a=function(e){return(null===e||void 0===e?void 0:e.clientHeight)||0},u=function(e){return e.classList.contains("sticky-shown")}},8629:function(e){e.exports={tooltip:"Tooltip_tooltip__DWDXo"}},94672:function(e){e.exports={container:"SavedSearchedInput_container__r57En",input:"SavedSearchedInput_input__YEKEd"}}}]);