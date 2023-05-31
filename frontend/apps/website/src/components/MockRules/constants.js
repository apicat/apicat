// 默认值

import { parseImage, parseStringWithType } from './parser'

export const mockSupportedLang = ['en', 'zh']
// 正则组合 校验语法
export const createDateTimeRegExp = (name) => new RegExp(`^${name}(\\|.{1,30})?$`)
export const createImageRegExp = (name) => new RegExp(`^${name}(\\|((\\d|[1-9]\\d+),(\\d|[1-9]\\d+)))?(\\|(\\d|[1-9]\\d+))?$`)
export const createOneOfRegExp = (name, types) => new RegExp(`^${name}(\\|(${types.join('|')}))?$`)
export const createOneOfWithRangeRegExp = (name, types) =>
  new RegExp(`^${name}(\\|(((${types.join('|')}),)?((\\d|[1-9]\\d+),(\\d|[1-9]\\d+)|(\\d|[1-9]\\d+))|(${types.join('|')})))?$`)
// export const createOneOfWithRangeRegExp = (name) => new RegExp(`^${name}(\\|((\\w+),)((\\d|[1-9]\\d+),(\\d|[1-9]\\d+)|(\\d|[1-9]\\d+)))?$`)
export const createStringRangeRegExp = (name) => new RegExp(`^${name}(\\((${mockSupportedLang.join('|')})\\))?(\\|((\\d|[1-9]\\d+),(\\d|[1-9]\\d+)|(\\d|[1-9]\\d+)))?$`)
export const createIntegerRangeRegExp = (name) => new RegExp(`^${name}(\\|\\-?((\\d|[1-9]\\d+),\\-?(\\d|[1-9]\\d+)|\\-?(\\d|[1-9]\\d+)))?$`)
export const creatFloatRangeRegExp = (name) => new RegExp(`^${name}(\\|((-?\\d+(\\.\\d+)?)|((-?\\d+(\\.\\d+)?),(-?\\d+(\\.\\d+)?),([1-9]\\d*))))?$`)
export const creatLangRegExp = (name) => new RegExp(`^${name}(\\((${mockSupportedLang.join('|')})\\))?$`)

/**
 * mock 语法列表
 */
const mockRules = {
  isUse: false,
  rules: {
    string: [
      {
        searchKey: '',
        name: 'string',
        cnName: '字符串',
        allow: {
          range: { min: 1, max: 10000 },
          // oneOfTypes: ['upper', 'letter', 'ansic', 'number'],
          regexp: createOneOfWithRangeRegExp('string', ['upper', 'letter', 'ansic', 'number']),
          parse: parseStringWithType,
        },
        syntax: [
          'string|{type?},{min?},{max?}',
          'type: number、upper、letter、ansic',
          'type为空时,默认为[a-zA-Z0-9]随机生成',
          'number 随机长度的字符串类型数字',
          'upper 随机长度的大写字符串',
          'letter 随机长度的小写字符串',
          'ansic 随机大小写字母数字以及一些特殊符号',
        ],
        example: [
          'string',
          '-> cegikmoq',
          'string|3',
          '-> wtq',
          'string|1,5',
          '-> xvt',
          'string|number',
          '-> 1234567890',
          'string|number,1,3',
          '-> 02',
          'string|upper',
          '-> GILXLDRTLX',
          'string|upper,1,3',
          '-> GI',
        ],
      },
      {
        searchKey: '',
        name: 'paragraph',
        cnName: '英文段落',
        allow: { actionText: '段落个数', range: { min: 1, max: 20 }, regexp: createStringRangeRegExp('paragraph') },
        syntax: [
          'paragraph(lang?)|{min?},{max?}',
          'lang: zh、en, 默认en',
          'paragraph',
          '随机生成大于等于3小于等于7个句子的英文段落',
          'paragraph|len',
          '生成指定句子个数的段落',
          'paragraph|min,max',
          '随机生成大于等于min小于等于max个句子的英文段落',
        ],
        example: [
          'paragraph',
          '-> Lfk zzz slewpibtmf tnhbu kufpzkuf fkp eimquyd vrnjfbwsok eimquy lwitfqcn zzzzzzzzz qhx upkfzu kufp vrnjfbw qhxofvmdt. Lgmsy lwit mylxkw lwitfqcny hovdkry lwitfqcn upkfzupkfz wtqnkhebxu xvtrpnl pfuk iqyhpxg eimquy aaaaaaaa aaaaaaa. Hcegikmo yxwvutsr odrgu aaaaaaa xvt yxwvutsr pfu qhxo vrnj qhxo nboc eimquydhl iqyhpxgow pfukzpfukz upkf xvtrpnljhf. Tkufp gmsyflr aaaaaa slewp odrguj vrn mylxk hov vrnjfbw zzzzzzzzzz dgjmpsv odrg eimq dgjmpsvy gmsyflr. Rupk rjbskc zzzzzzz gmsyfl nbocpdqer slewpibt zzzzzzzzzz jscluenwg zzzzz bcdefgh xvtrp wtqnkheb wtqnk. Lmylxkwjv wtqnkh bcde cegikm tnhb gmsyflrx odrgujxm aaaa iqyhpx sle iqyhp upkfzupkf pfukzp upkfzupk iqy aaaaa nbo dgj. Zrjbskc hov nbo mylxkwjviu bcde wtqnkhe aaaaaa jscl lwitfqc upk lwitfqcnyk xvtrpnljhf yxwvutsr fkpuzf.',
          'paragraph|2',
          '-> Wkufpzkufp kufpz lwi slewpibt tnh zzzzzz wtqnkhe odrgujxmb yxwvutsrq tnhbuoicv nboc zzz odrgujxmbp zzzzzzzzzz jsclue. Rcegikm yxwvutsrq aaaaaaa hovd xvt qhxofvmdtk wtqnkheb slewpibtm mylx vrnjfbw qhxofvmdtk pfukzpfukz wtqn.',
          'paragraph|1-3',
          '-> Lxvtrpnl vrn qhxofvmdt lwitfqcny yxwv tnhbuo iqyhpxgowf rjbskc jscluenw yxwvutsrq fkpuzfkp upkfzup vrnjfbws upkfzu aaaaaaa lwit tnhbuoicv. Fzzz eimq aaaaa iqyh mylxkwjv vrnjfbwsok iqyhpxgo gmsyfl lwitfqcnyk slewpibtm upkfzupk kufpzk upkfzupkfz slewpibtmf yxwvutsrqp qhxofvmd nbocpdqe.',
        ],
      },
      {
        searchKey: '',
        name: 'sentence',
        cnName: '英文句子',
        allow: { actionText: '单词个数', range: { min: 1, max: 30 }, regexp: createStringRangeRegExp('sentence') },
        syntax: [
          'sentence',
          '随机生成大于等于12小于等于18个单词的英文句子',
          'sentence|len',
          '生成指定单词个数的句子',
          'sentence|min,max',
          '随机生成大于等于min小于等于max个单词的英文句子',
        ],
        example: [
          'sentence',
          '-> Okufpz hovdkr fkpuzf pfukzpfukz aaaaaaaa gmsyf odrgu fkpuzfkp lwitf rjbs eim nbocpdqer mylx hovdk.',
          'sentence|5',
          '-> Mjsclu xvtrpn slewp xvtrpn slewp.',
          'sentence|3,5',
          '-> Epfukzpf lwitfqcn upkfzupk.',
        ],
      },
      {
        searchKey: '',
        name: 'word',
        cnName: '英文单词',
        allow: { actionText: '字母个数', range: { min: 1, max: 20 }, regexp: createStringRangeRegExp('word') },
        syntax: ['word', '随机生成大于等于3小于等于10个字母的英文单词', 'word|len', '生成指定字母个数的单词', 'word|min,max', '随机生成大于等于min小于等于max个字母的英文单词'],
        example: ['word', '-> aaa', 'word|5', '-> slewp', 'word|3,5', '-> yxwvu'],
      },
      {
        searchKey: '',
        name: 'title',
        cnName: '英文标题',
        allow: { actionText: '单词个数', range: { min: 1, max: 15 }, regexp: createStringRangeRegExp('title') },
        syntax: [
          'title',
          '随机生成大于等于3小于等于7个单词的英文标题',
          'title|len',
          '生成指定单词个数的英文标题',
          'title|min,max',
          '随机生成大于等于min小于等于max个单词的英文标题',
        ],
        example: ['title', '-> Omylxk Ubcdefghij Cfkpuz Kcegikm', 'title|5', '-> Bslewpi Tbcde Xslewpibtm Thovdkrygn Ymylxk', 'title|3,5', '-> Pgmsyflr Qwtqnkheb Sup Yxvtrpnljh'],
      },
      {
        searchKey: '',
        name: 'phrase',
        cnName: '短语',
        allow: { actionText: '短语个数', range: { min: 1, max: 15 }, regexp: createStringRangeRegExp('phrase') },
        syntax: [
          'phrase(lang?)|{min?},{max?}',
          'lang: zh、en, 默认en',
          'phrase',
          '随机生成大于等于3小于等于7个的英文短语',
          'phrase|len',
          '生成指定个数的英文短语',
          'phrase|min,max',
          '随机生成大于等于min小于等于max个的英文短语',
        ],
        example: [
          'phrase',
          '-> Omylxk Ubcdefghij Cfkpuz Kcegikm',
          'phrase|5',
          '-> Bslewpi Tbcde Xslewpibtm Thovdkrygn Ymylxk',
          'phrase|3,5',
          '-> Pgmsyflr Qwtqnkheb Sup Yxvtrpnljh',
        ],
      },
      {
        searchKey: '',
        name: 'phone',
        cnName: '手机号',
        allow: { regexp: creatLangRegExp('phone') },
        syntax: ['phone(lang?)', 'lang: zh、en, 默认en', 'phone', '随机生成一个手机号码'],
        example: ['phone(zh)', '-> 13333333333'],
      },

      {
        searchKey: '',
        name: 'idcard',
        cnName: '身份证号',
        allow: { regexp: creatLangRegExp('idcard') },
        syntax: ['idcard(lang?)', 'lang: zh、en, 默认en', 'idcard', '随机生成一个身份证号码'],
        example: ['idcard', '-> 610102202003170019'],
      },
      {
        searchKey: '',
        name: 'url',
        cnName: '网址',
        syntax: ['url', '随机生成一个url'],
        example: ['url', '-> https://apicat.net'],
      },
      {
        searchKey: '',
        name: 'domain',
        cnName: '域名',
        syntax: ['domain', '随机生成一个域名'],
        example: ['domain', '-> apicat.net'],
      },
      {
        searchKey: 'ip ipv6',
        name: 'ipv6',
        cnName: 'IPv6地址',
        syntax: ['ipv6', '随机生成一个ipv6地址'],
        example: ['ipv6', '-> 2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b'],
      },
      {
        searchKey: 'ip ipv4',
        name: 'ipv4',
        cnName: 'IPv4地址',
        syntax: ['ipv4', '随机生成一个ipv4地址'],
        example: ['ipv4', '-> 127.0.0.1'],
      },
      {
        searchKey: '',
        name: 'email',
        cnName: '电子邮箱',
        syntax: ['email', '随机生成一个email'],
        example: ['email', '-> helloworld@gmail.com'],
      },
      {
        searchKey: '省份',
        name: 'provinceorstate',
        cnName: '省',
        allow: { regexp: creatLangRegExp('provinceorstate') },
        syntax: ['provinceorstate(lang?)', 'lang: zh、en, 默认en', 'provinceorstate', '随机生成一个省名称'],
        example: ['provinceorstate(zh)', '-> 陕西省'],
      },
      {
        searchKey: '省份 城市',
        name: 'provinceorstatecity',
        cnName: '省市',
        allow: { regexp: creatLangRegExp('provinceorstatecity') },
        syntax: ['provinceorstatecity(lang?)', 'lang: zh、en, 默认en', 'provinceorstatecity', '随机生成一个省市名称'],
        example: ['provinceorstatecity(zh)', '-> 陕西省 西安市'],
      },
      {
        searchKey: '',
        name: 'city',
        cnName: '城市',
        allow: { regexp: creatLangRegExp('city') },
        syntax: ['city(lang?)', 'lang: zh、en, 默认en', 'city', '随机生成一个城市名称'],
        example: ['city(zh)', '-> 西安'],
      },
      {
        searchKey: '',
        name: 'street',
        cnName: '街道',
        allow: { regexp: creatLangRegExp('street') },
        syntax: ['street(lang?)', 'lang: zh、en, 默认en', 'street', '随机生成街道名称'],
        example: ['street(zh)', '-> 因看大街吴云 商场64-4号'],
      },
      {
        searchKey: '地址',
        name: 'address',
        cnName: '省市区',
        allow: { regexp: creatLangRegExp('address') },
        syntax: ['address(lang?)', 'lang: zh、en, 默认en', 'address', '随机生成一个详细地址'],
        example: ['address(zh)', '-> 陕西省西安市新城区还用场集小区1栋0单元975号'],
      },
      {
        searchKey: '',
        name: 'zipcode',
        cnName: '邮政编码',
        allow: { regexp: creatLangRegExp('zipcode') },
        syntax: ['zipcode(lang?)', 'lang: zh、en, 默认en', 'zipcode', '随机生成一个邮政编码'],
        example: ['zipcode', '-> 710000'],
      },
      {
        searchKey: '',
        name: 'longitude',
        cnName: '经度',
        syntax: ['longitude', '随机生成一个经度'],
        example: ['longitude', '-> 116.397128'],
      },
      {
        searchKey: '',
        name: 'latitude',
        cnName: '维度',
        syntax: ['latitude', '随机生成一个维度'],
        example: ['latitude', '-> 39.916527'],
      },
      {
        searchKey: 'longitude latitude',
        name: 'longitudelatitude',
        cnName: '经纬度',
        syntax: ['longitudelatitude', '随机生成一个经纬度'],
        example: ['longitudelatitude', '-> 116.397128, 39.916527'],
      },
      {
        searchKey: '',
        name: 'date',
        cnName: '日期',
        allow: { regexp: createDateTimeRegExp('date') },
        syntax: ['date|{format?}'],
        example: ['date', '-> 2020-03-17', 'date|YYYY年MM月dd日', '-> 2020年01月20日', 'date|YYYY-MM-dd', '-> 2020-03-17', 'date|YYYY/MM/dd', '-> 2020/03/17'],
      },
      {
        searchKey: '',
        name: 'time',
        cnName: '时间',
        allow: { regexp: createDateTimeRegExp('time') },
        syntax: ['time', '随机生成一个(HH:mm:ss)格式的时间', 'time|format', '随机生成一个指定格式的时间'],
        example: ['time', '-> 09:00:00', 'time|HH:mm:ss', '-> 09:00:00'],
      },
      {
        searchKey: '',
        name: 'datetime',
        cnName: '日期时间',
        allow: { regexp: createDateTimeRegExp('datetime') },
        syntax: ['datetime|{format?}'],
        example: [
          'datetime',
          '-> 2006-01-02T15:04:05Z07:00',
          'datetime|"YYYY年MM月dd日 HH:mm"',
          '-> 2020年01月20日 12:00',
          'date|YYYY-MM-dd',
          '-> 2020-03-17',
          'date|YYYY/MM/dd',
          '-> 2020/03/17',
        ],
      },
      {
        searchKey: '',
        name: 'timestamp',
        cnName: '时间戳',
        syntax: ['timestamp', '随机生成一个时间戳'],
        example: ['timestamp', '-> 1584406800'],
      },
      {
        searchKey: '',
        name: 'now',
        allow: { regexp: createDateTimeRegExp('now') },
        cnName: '当前时间',
        syntax: ['now', '当前时间'],
        example: ['now', '-> 2022-05-22 13:00:32'],
      },
      {
        searchKey: '',
        name: 'imagedata',
        cnName: '图片数据',
        allow: {
          isSwap: true,
          range: { min: 20, max: 1024, minActionText: '图片宽度', maxActionText: '图片高度' },
          parse: parseImage,
          regexp: createImageRegExp('imagedata'),
        },
        syntax: [
          'imagedata|{width?},{height?}',
          '生成一个128*128的Base64编码图片',
          'imagedata|width',
          '生成一个宽高相等的Base64编码图片',
          'imagedata|width,height',
          '生成一个指定宽和高的Base64编码图片',
        ],
        example: [
          'imagedata',
          '-> data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAM...',
          'imagedata|200',
          '-> data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAM...',
          'imagedata|600*400',
          '-> data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAM...',
        ],
      },
      {
        searchKey: 'image url',
        name: 'imageurl',
        cnName: '图片链接',
        allow: {
          isSwap: true,
          range: { min: 20, max: 1024, minActionText: '图片宽度', maxActionText: '图片高度' },
          parse: parseImage,
          regexp: createImageRegExp('imageurl'),
        },
        syntax: [
          'imageurl|{width?},{height?}',
          '生成一个128*128的图片链接地址',
          'imageurl|width',
          '生成一个宽高相等的图片链接地址',
          'imageurl|width,height',
          '生成一个指定宽和高的图片链接地址',
        ],
        example: [
          'imageurl',
          '-> https://dummyimage.com/128/128.png',
          'imageurl|200',
          '-> https://dummyimage.com/200/200.png',
          'imageurl|200,100',
          '-> https://dummyimage.com/200/100.png',
        ],
      },
      {
        searchKey: '',
        name: 'color',
        cnName: '颜色',
        allow: { regexp: createOneOfRegExp('color', ['rgb', 'rgba', 'hsl', 'hex']) },
        syntax: ['color', '默认生成hex颜色', 'color|type', '随机生成颜色值,类型type只支持rgb、rgba、hsl、hex'],
        example: ['color', '-> #002211', 'color|rgb', '-> rgb(255,0,123)'],
      },
      {
        searchKey: '状态码',
        name: 'httpcode',
        cnName: 'http状态码',
        syntax: ['httpcode', '随机生成一个http状态码'],
        example: ['httpcode', '-> 200'],
      },
      {
        searchKey: 'http 请求方法',
        name: 'httpmethod',
        cnName: 'http请求方法',
        syntax: ['httpmethod', '随机生成一个http请求方法'],
        example: ['httpmethod', '-> DELETE'],
      },
      {
        searchKey: 'uuid id',
        name: 'uuid',
        cnName: '唯一ID',
        syntax: ['uuid', '随机生成一个uuid字符串'],
        example: ['uuid', '-> 9b271dc8-abb9-19b0-f5f4-43225ff7968c'],
      },
      {
        searchKey: '数字占位符 占位符',
        name: 'numberpattern',
        defaultValue: 'numberpattern|###-######-#',
        cnName: '数字占位',
        allow: { regexp: /^numberpattern\|[#-]+$/ },
        syntax: ['numberpattern|{pattern}', 'pattern为#和-组合的方式'],
        example: ['numberpattern|###-######-#', '-> 111-222222-3'],
      },
      {
        searchKey: '名字 姓 名',
        name: 'firstname',
        cnName: '姓',
        allow: { regexp: creatLangRegExp('firstname') },
        syntax: ['firstname(lang?)', 'lang: zh、en, 默认en', '随机生成一个姓'],
        example: ['firstname', '-> Robert', 'firstname(zh)', '-> 王'],
      },
      {
        searchKey: '名字 姓 名',
        name: 'name',
        cnName: '名字',
        allow: { regexp: creatLangRegExp('name') },
        syntax: ['name(lang?)', 'lang: zh、en, 默认en', 'name', '随机生成一个名字'],
        example: ['name', '-> Robert Robinson', 'name(zh)', '-> 王发财'],
      },
      {
        searchKey: '名字 姓 名',
        name: 'lastname',
        cnName: '名',
        allow: { regexp: creatLangRegExp('lastname') },
        syntax: ['lastname(lang?)', 'lang: zh、en, 默认en', 'lastname', '随机生成一个名'],
        example: ['lastname', '-> Robinson', 'lastname(zh)', '-> 发财'],
      },
      {
        searchKey: '性别',
        name: 'gender',
        cnName: '性别',
        allow: { regexp: creatLangRegExp('gender') },
        syntax: ['gender(lang?)', 'lang: zh、en, 默认en', 'gender', '随机生成性别'],
        example: ['gender(zh)', '-> 男', 'gender', '-> female'],
      },
      {
        searchKey: 'one of 单选',
        name: 'oneof',
        defaultValue: 'oneof|0,1',
        cnName: '单选',
        allow: { regexp: /^oneof\|([^,]+(,[^,]+)*)$/ },
        syntax: ['oneof|{value?}...', 'value1,value2中随机选择一个。语法:以英文逗号","分割数据。'],
        example: ['oneof|男,女', '-> 女', 'oneof|1,"a b",a,b', '-> a'],
      },
    ],
    integer: [
      {
        searchKey: '',
        name: 'integer',
        cnName: '整数',
        allow: { range: { min: -1000000, max: 1000000 }, regexp: createIntegerRangeRegExp('integer') },
        syntax: ['integer', '随机生成0-1000的整数', 'integer|count', '生成指定数值的整数', 'integer|min,max', '随机生成大于等于min小于等于max的整数'],
        example: ['integer', '-> 123', 'integer|123', '-> 123', 'integer|1,10', '-> 5'],
      },
      {
        searchKey: '字增 ID',
        name: 'autoincrement',
        cnName: '自增',
        allow: { isSwap: true, range: { min: -1000000, max: 1000000 }, regexp: createIntegerRangeRegExp('autoincrement') },
        syntax: ['autoincrement|{begin?},{step?}', 'begin起始值,默认 1', 'step步长,默认 1'],
        example: ['autoincrement', '-> 1,2,3,....', 'autoincrement|100', '-> 100,101,102,...', 'autoincrement|100,2', '-> 100,102,104,...'],
      },
      {
        searchKey: '',
        name: 'phone',
        cnName: '手机号',
        allow: { regexp: creatLangRegExp('phone') },
        syntax: ['phone(lang?)', 'lang: zh、en, 默认en', 'phone', '随机生成一个手机号码'],
        example: ['phone', '-> 13333333333'],
      },
      {
        searchKey: '',
        name: 'idcard',
        cnName: '身份证号',
        allow: { regexp: creatLangRegExp('idcard') },
        syntax: ['idcard(lang?)', 'lang: zh、en, 默认en', 'idcard', '随机生成一个身份证号码'],
        example: ['idcard', '-> 610102202003170019'],
      },
      {
        searchKey: '',
        name: 'zipcode',
        cnName: '邮政编码',
        allow: { regexp: creatLangRegExp('zipcode') },
        syntax: ['zipcode', '随机生成一个邮政编码'],
        example: ['zipcode', '-> 710000'],
      },
      {
        searchKey: '',
        name: 'timestamp',
        cnName: '时间戳',
        syntax: ['timestamp', '随机生成一个时间戳'],
        example: ['timestamp', '-> 1584406800'],
      },
    ],
    boolean: [
      {
        searchKey: '',
        name: 'boolean',
        cnName: '布尔值',
        allow: {
          // range: { min: 0, max: 100 },
          // parse: parseBoolean,
          regexp: createOneOfRegExp('boolean', ['true', 'false']),
        },
        syntax: ['boolean', '随机生成一个50%概率的true or false', 'boolean|value', 'value为true or false,代表生成指定的布尔值'],
        example: ['boolean', '-> true', 'boolean|true', '-> true', 'boolean|false', '-> false'],
      },
    ],
    array: [
      {
        searchKey: '',
        name: 'array',
        cnName: '数组',
        syntax: ['array', '数组内的所有元素,随机循环次数'],
        example: ['// 数组arr1', 'arr1 = [string|1]', '// 数组arr1的规则及结果', 'array', '-> ["a", "b", "c"]'],
      },
    ],
    object: [
      {
        searchKey: '',
        name: 'object',
        cnName: '对象',
        syntax: ['object', '对象内的所有元素以key:val的方式组成object'],
        example: ['// 对象obj', 'obj = {"username": string|6, "password": string|6}', '// 对象obj的规则及结果', 'object', '-> {"username": "abcdef", "password": "123456"}'],
      },
    ],
    file: [
      {
        searchKey: '',
        name: 'file',
        cnName: '文件',
        allow: { regexp: createOneOfRegExp('file', ['word', 'excel', 'csv', 'md']) },
        syntax: ['file', '生成一个Markdown文件,以文件流的方式返回', 'file|type', '生成一个指定类型的文件,以文件流的方式返回,类型type只支持word、excel、csv、md'],
        example: ['file', '-> 返回一个Markdown文件流内容', 'file|word', '-> 返回一个word文件流'],
      },
      {
        searchKey: '',
        name: 'image',
        cnName: '图片',
        allow: {
          isSwap: true,
          range: { min: 20, max: 1024, minActionText: '图片宽度', maxActionText: '图片高度' },
          parse: parseImage,
          regexp: createImageRegExp('image'),
        },
        syntax: [
          'image',
          '生成一个200*150的jpeg图片流',
          'image|width,height',
          '随机生成一个指定宽高的jpeg图片流',
          'image|width,height',
          '随机生成一个指定宽高和类型的图片流,类型type只支持jpeg和png',
        ],
        example: ['image', '-> 返回一个200*150的jpeg图片流', 'image|200', '-> 返回一个200*200的jpeg图片流', 'image|600*400,png', '-> 返回一个600*400的png图片流'],
      },
    ],
    number: [
      {
        searchKey: '',
        name: 'float',
        cnName: '浮点数',
        allow: {
          prefix: '整数',
          actionText: '整数',
          range: { min: -1000000, max: 1000000 },
          regexp: creatFloatRangeRegExp('float'),
        },
        syntax: [
          'float|{min?},{max?},{fixed?}',
          'float',
          '随机一个浮点数',
          'float|浮点数',
          '生成指定的浮点数',
          'float|min,max,fixed',
          '随机生成大于等于min小于等于max的整数,小数部分保留fixed位的随机小数',
        ],
        example: ['float', '// 随机一个浮点数', 'float|102.01', '-> 102.01', 'float|10000,20000,4', '-> 13886.1021'],
      },
    ],
  },
}

const rules = {}

export const getMockRules = () => {
  if (mockRules.isUse) {
    return rules
  }

  Object.keys(mockRules.rules).forEach((key) => {
    const _rules = mockRules.rules[key] || []

    const obj = {
      ruleKeys: [],
      rules: [],
    }

    obj.rules = _rules.map((item) => {
      obj.ruleKeys.push(item.name)

      item.syntax = item.syntax.join('<br/>')
      item.example = item.example.join('<br/>')
      item.searchKey = item.searchKey.split(' ').concat([item.name, item.cnName, key]).join(' ')
      item.searchKeys = item.searchKey.split(' ').concat([item.name, item.cnName, key])
      item.key = key + '-' + item.name
      return item
    })

    rules[key] = obj
  })

  mockRules.isUse = true

  return rules
}
