import iconBlockH1 from '../../assets/images/icon-block-h1.png'
import iconBlockH2 from '../../assets/images/icon-block-h2.png'
import iconBlockH3 from '../../assets/images/icon-block-h3.png'
import iconBlockImg from '../../assets/images/icon-block-img.png'
import iconBlockBulletlist from '../../assets/images/icon-block-bulletlist.png'
import iconBlockOrderedllist from '../../assets/images/icon-block-orderedllist.png'
import iconBlockBlockquote from '../../assets/images/icon-block-blockquote.png'
import iconBlockHorizontalrule from '../../assets/images/icon-block-horizontalrule.png'
import iconBlockLink from '../../assets/images/icon-block-link.png'
import iconBlockTable from '../../assets/images/icon-block-table.png'
import iconBlockCodeblock from '../../assets/images/icon-block-codeblock.png'
import iconBlockHttpapiurl from '../../assets/images/icon-block-httpapiurl.png'
import iconBlockHttpstatuscode from '../../assets/images/icon-block-httpstatuscode.png'
import iconBlockApirequest from '../../assets/images/icon-block-apirequest.png'
import iconBlockApiresponse from '../../assets/images/icon-block-apiresponse.png'
import iconBlockApiurl from '../../assets/images/icon-block-apiurl.png'
import iconBlockApiparameter from '../../assets/images/icon-block-apiparameter.png'
import iconBlockNotification from '../../assets/images/icon-block-notification.png'
import { isMac } from '../utils'

const mod = isMac ? '⌘' : 'ctrl'

export default function blockMenuItems() {
    return [
        {
            name: 'separator',
            title: '基础模块',
        },

        {
            name: 'heading',
            title: '一级标题 (Heading 1)',
            keywords: 'h1 heading1 title',
            icon: 'icon-block-h1',
            img: iconBlockH1,
            shortcut: '^ ⇧ 1',
            desc: '创建一个大节标题',
            attrs: { level: 1 },
        },
        {
            name: 'heading',
            title: '二级标题 (Heading 2)',
            keywords: 'h2 heading2 title',
            icon: 'icon-block-h2',
            img: iconBlockH2,
            shortcut: '^ ⇧ 2',
            desc: '创建一个中段标题',
            attrs: { level: 2 },
        },
        {
            name: 'heading',
            title: '三级标题 (Heading 3)',
            keywords: 'h3 heading3 title',
            icon: 'icon-block-h3',
            img: iconBlockH3,
            shortcut: '^ ⇧ 3',
            desc: '创建一个小节标题',
            attrs: { level: 3 },
        },
        {
            name: 'bullet_list',
            title: '无序列表 (Bullet list)',
            icon: 'icon-block-bulletlist',
            img: iconBlockBulletlist,
            shortcut: '^ ⇧ 8',
            keywords: 'bullet list',
            desc: '创建一个无序列表',
        },
        {
            name: 'ordered_list',
            title: '有序列表 (Ordered list)',
            icon: 'icon-block-orderedllist',
            img: iconBlockOrderedllist,
            shortcut: '^ ⇧ 9',
            keywords: 'ordered list',
            desc: '创建一个有序列表',
        },
        {
            name: 'blockquote',
            title: '块引用 (Blockquote)',
            icon: 'icon-block-blockquote',
            img: iconBlockBlockquote,
            keywords: 'block quote',
            shortcut: `${mod} ]`,
            desc: '创建一个块引用',
        },
        {
            name: 'hr',
            title: '分割线 (Horizontal rule)',
            icon: 'icon-block-horizontalrule',
            img: iconBlockHorizontalrule,
            shortcut: `${mod} _`,
            keywords: 'horizontal rule break line',
            desc: '创建一个水平分割线',
        },
        {
            name: 'link',
            title: '超链接 (Link)',
            icon: 'icon-block-link',
            img: iconBlockLink,
            shortcut: `${mod} k`,
            keywords: 'link url uri href',
            desc: '创建一个超链接',
        },
        {
            name: 'table',
            title: '表格 (Table)',
            icon: 'icon-block-table',
            img: iconBlockTable,
            keywords: 'table',
            attrs: { rowsCount: 3, colsCount: 3 },
            desc: '创建一个表格',
        },
        {
            name: 'image',
            title: '图片 (Image)',
            img: iconBlockImg,
            keywords: 'picture photo image img',
            desc: '上传一个图片',
        },
        {
            name: 'code_block',
            title: '代码块 (Code block)',
            icon: 'icon-block-codeblock',
            img: iconBlockCodeblock,
            shortcut: '^ ⇧ \\',
            keywords: 'code',
            desc: '创建一个代码块',
        },

        {
            name: 'separator',
            title: '接口模块',
        },

        {
            name: 'http_api_url',
            title: 'HTTP接口访问地址 (HTTP API URL)',
            keywords: 'http url',
            icon: 'icon-block-httpapiurl',
            img: iconBlockHttpapiurl,
            desc: '创建一个基于HTTP协议的接口请求地址',
        },
        {
            name: 'http_status_code',
            title: 'HTTP状态码 (HTTP Status Code)',
            keywords: 'http code',
            icon: 'icon-block-httpstatuscode',
            img: iconBlockHttpstatuscode,
            desc: '创建一个HTTP状态码',
        },
        {
            name: 'http_api_request_parameter',
            title: 'HTTP接口请求参数 (HTTP API Request Parameter)',
            keywords: 'http request param table',
            icon: 'icon-block-apirequest',
            img: iconBlockApirequest,
            desc: '创建一个基于HTTP协议的接口请求参数表格',
        },
        {
            name: 'http_api_response_parameter',
            title: 'HTTP接口响应参数 (HTTP API Response Parameter)',
            keywords: 'http response param table',
            icon: 'icon-block-apiresponse',
            img: iconBlockApiresponse,
            desc: '创建一个基于HTTP协议的接口响应参数表格',
        },
        {
            name: 'api_url',
            title: '接口访问地址 (API URL)',
            keywords: 'api url',
            icon: 'icon-block-apiurl',
            img: iconBlockApiurl,
            desc: '创建一个接口请求地址',
        },
        {
            name: 'api_parameter',
            title: '接口参数 (API Parameter)',
            keywords: 'param table api',
            icon: 'icon-block-apiparameter',
            img: iconBlockApiparameter,
            desc: '创建一个接口参数表格',
        },

        {
            name: 'separator',
            title: '团队协作',
        },

        // {
        //     name: 'trigger_openNotification',
        //     title: '提醒成员(@)',
        //     keywords: '@',
        //     icon: 'icon-block-notification',
        //     img: iconBlockNotification,
        //     desc: '提醒成员关注文档变化',
        // },
    ]
}
