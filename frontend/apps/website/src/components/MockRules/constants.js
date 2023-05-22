// 默认值

import { parseImage, parseBoolean } from './parser'

// form data 基础类型，mock数据类型 1.int 2.float 3.string 4.array 5.object 6.boolean 7.file
export const PARAM_TYPES = {
  TYPES: [
    { text: 'Int', value: 1 },
    { text: 'Float', value: 2 },
    { text: 'String', value: 3 },
    { text: 'Array', value: 4 },
    { text: 'Object', value: 5 },
    { text: 'ArrayObject', value: 8 },
    { text: 'Boolean', value: 6 },
    { text: 'File', value: 7 },
  ],
  valueOf(value) {
    return (this.TYPES.find((item) => item.value === value) || { text: '' }).text
  },
  VALUES: {
    ARRAY: 4,
  },
}

// 正整数
export const RE_STR_NORMAL_NUMBER = '(\\-?(\\d|[1-9]\\d+))'

export const createDateTimeRegExp = (name) => new RegExp(`^${name}((\\|(y-m-d|y\\/m\\/d)) h:i:s)?$`, 'i')
export const createDateRegExp = (name) => new RegExp(`^${name}(\\|(y-m-d|y\\/m\\/d))?$`, 'i')
export const createTimeRegExp = (name) => new RegExp(`^${name}(\\|h:i:s)?$`, 'i')
export const createImageRegExp = (name) => new RegExp(`^${name}(\\|((\\d|[1-9]\\d+)\\*(\\d|[1-9]\\d+)),(\\w+))?(\\|((\\d|[1-9]\\d+)\\*(\\d|[1-9]\\d+)))?$`)
export const createOneOfRegExp = (name, types) => new RegExp(`^${name}(\\|(${types.join('|')}))?$`)
export const createRangeRegExp = (name) => new RegExp(`^(${name})(\\|((\\-?(\\d|[1-9]\\d+))\\-(\\d|[1-9]\\d+)|(\\-?(\\d|[1-9]\\d+))))?$`)
export const createNumberRangeRegExp = (name) => new RegExp(`^${name}(\\|\\-?((\\d|[1-9]\\d+)~\\-?(\\d|[1-9]\\d+)|\\-?(\\d|[1-9]\\d+)))?$`)
export const creatFloatRangeRegExp = (name) => new RegExp(`^${name}(\\|\\-?((\\d|[1-9]\\d+)~\\-?(\\d|[1-9]\\d+)|\\-?(\\d|[1-9]\\d+))(?:\\.(\\d+-?\\d*)))?$`)

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
        allow: { range: { min: 1, max: 10000 }, regexp: createRangeRegExp('string') },
        syntax: [
          'string',
          '随机生成大于等于3小于等于10长度的英文字符串',
          'string|len',
          '生成指定长度的英文字符串',
          'string|min-max',
          '随机生成大于等于min小于等于max长度的英文字符串',
        ],
        example: ['string', '-> cegikmoq', 'string|3', '-> wtq', 'string|1-5', '-> xvt'],
      },
      {
        searchKey: '',
        name: 'paragraph',
        cnName: '英文段落',
        allow: { actionText: '句子个数', range: { min: 1, max: 20 }, regexp: createRangeRegExp('paragraph') },
        syntax: [
          'paragraph',
          '随机生成大于等于3小于等于7个句子的英文段落',
          'paragraph|len',
          '生成指定句子个数的段落',
          'paragraph|min-max',
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
        allow: { actionText: '单词个数', range: { min: 1, max: 30 }, regexp: createRangeRegExp('sentence') },
        syntax: [
          'sentence',
          '随机生成大于等于12小于等于18个单词的英文句子',
          'sentence|len',
          '生成指定单词个数的句子',
          'sentence|min-max',
          '随机生成大于等于min小于等于max个单词的英文句子',
        ],
        example: [
          'sentence',
          '-> Okufpz hovdkr fkpuzf pfukzpfukz aaaaaaaa gmsyf odrgu fkpuzfkp lwitf rjbs eim nbocpdqer mylx hovdk.',
          'sentence|5',
          '-> Mjsclu xvtrpn slewp xvtrpn slewp.',
          'sentence|3-5',
          '-> Epfukzpf lwitfqcn upkfzupk.',
        ],
      },
      {
        searchKey: '',
        name: 'word',
        cnName: '英文单词',
        allow: { actionText: '字母个数', range: { min: 1, max: 20 }, regexp: createRangeRegExp('word') },
        syntax: ['word', '随机生成大于等于3小于等于10个字母的英文单词', 'word|len', '生成指定字母个数的单词', 'word|min-max', '随机生成大于等于min小于等于max个字母的英文单词'],
        example: ['word', '-> aaa', 'word|5', '-> slewp', 'word|3-5', '-> yxwvu'],
      },
      {
        searchKey: '',
        name: 'title',
        cnName: '英文标题',
        allow: { actionText: '单词个数', range: { min: 1, max: 15 }, regexp: createRangeRegExp('title') },
        syntax: [
          'title',
          '随机生成大于等于3小于等于7个单词的英文标题',
          'title|len',
          '生成指定单词个数的英文标题',
          'title|min-max',
          '随机生成大于等于min小于等于max个单词的英文标题',
        ],
        example: ['title', '-> Omylxk Ubcdefghij Cfkpuz Kcegikm', 'title|5', '-> Bslewpi Tbcde Xslewpibtm Thovdkrygn Ymylxk', 'title|3-5', '-> Pgmsyflr Qwtqnkheb Sup Yxvtrpnljh'],
      },
      {
        searchKey: '',
        name: 'cparagraph',
        cnName: '中文段落',
        allow: { actionText: '句子个数', range: { min: 1, max: 15 }, regexp: createRangeRegExp('cparagraph') },
        syntax: [
          'cparagraph',
          '随机生成大于等于3小于等于7个句子的中文段落',
          'cparagraph|len',
          '生成指定句子个数的中文段落',
          'cparagraph|min-max',
          '随机生成大于等于min小于等于max个句子的中文段落',
        ],
        example: [
          'cparagraph',
          '-> 燥卜簇嬉藏扯逻梅久狡双红宦蛤拉酣概些。捏腰勋橱谈失猿阻怨精你昆蓖呜。臣毅娄韩缸某均铭等斟吹辨码挣勘豁瓶意。袒瞻朵害侵赚箭歼弄割酌聚湖。奶曹礼贝做币毅绣邮卦慎澎芋。北妆惦期漂禾舷费鞋兑器床摊。',
          'cparagraph|2',
          '-> 蓉熟战半进箍柬娘齐蒲熊成十近筷柔。肴栅喷述看所做蔫泼媳隘篷擎卦谆片度。',
          'cparagraph|1-3',
          '-> 监妖认机传秧尽赞榛再纠弓逝泼千肆愚镰。恭砖闽尔牌躁坟沽袄厢杜舰偏换约。',
        ],
      },
      {
        searchKey: '',
        name: 'csentence',
        cnName: '中文句子',
        allow: { actionText: '汉字个数', range: { min: 1, max: 30 }, regexp: createRangeRegExp('csentence') },
        syntax: [
          'csentence',
          '随机生成大于等于12小于等于18个汉字的中文句子',
          'csentence|len',
          '生成指定汉字个数的中文句子',
          'csentence|min-max',
          '随机生成大于等于min小于等于max个汉字的中文句子',
        ],
        example: ['csentence', '-> 袋瞪未审供赂答欧庞到遥翻淑。', 'csentence|5', '-> 夏斧皇衡中。', 'csentence|3-5', '-> 努壶心敷泣。'],
      },
      {
        searchKey: '',
        name: 'cword',
        cnName: '中文词语',
        allow: { actionText: '汉字个数', range: { min: 1, max: 10 }, regexp: createRangeRegExp('cword') },
        syntax: [
          'cword',
          '随机生成大于等于1小于等于4个汉字的中文词语',
          'cword|len',
          '生成指定汉字个数的中文词语',
          'cword|min-max',
          '随机生成大于等于min小于等于max个汉字的中文词语',
        ],
        example: ['cword', '-> 闪', 'cword|3', '-> 批纱人', 'cword|3-5', '-> 撤蓬喻爆隆'],
      },
      {
        searchKey: '',
        name: 'ctitle',
        cnName: '中文标题',
        allow: { actionText: '汉字个数', range: { min: 1, max: 20 }, regexp: createRangeRegExp('ctitle') },
        syntax: [
          'ctitle',
          '随机生成大于等于3小于等于7个汉字的中文标题',
          'ctitle|len',
          '生成指定汉字个数的中文标题',
          'ctitle|min-max',
          '随机生成大于等于min小于等于max个汉字的中文标题',
        ],
        example: ['ctitle', '-> 徊疚述坷毯蔫', 'ctitle|3', '-> 滤儿疯', 'ctitle|3-5', '-> 廊狼贿'],
      },

      {
        searchKey: '',
        name: 'mobile',
        cnName: '手机号',
        syntax: ['mobile', '随机生成一个国内手机号码'],
        example: ['mobile', '-> 13333333333'],
      },
      {
        searchKey: '',
        name: 'phone',
        cnName: '座机号',
        syntax: ['phone', '随机生成一个国内座机号码'],
        example: ['phone', '-> (029)88888888'],
      },
      {
        searchKey: '',
        name: 'idcard',
        cnName: '身份证号',
        syntax: ['idcard', '随机生成一个身份证号码'],
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
        searchKey: '',
        name: 'ip',
        cnName: 'IP地址',
        syntax: ['ip', '随机生成一个ip地址'],
        example: ['ip', '-> 127.0.0.1'],
      },
      {
        searchKey: '',
        name: 'email',
        cnName: '电子邮箱',
        syntax: ['email', '随机生成一个email'],
        example: ['email', '-> helloworld@gmail.com'],
      },
      {
        searchKey: '',
        name: 'province',
        cnName: '省份',
        syntax: ['province', '随机生成一个省份名称'],
        example: ['province', '-> 陕西省'],
      },
      {
        searchKey: '',
        name: 'city',
        cnName: '城市',
        syntax: ['city', '随机生成一个城市名称'],
        example: ['city', '-> 西安'],
      },
      {
        searchKey: '',
        name: 'province_city',
        cnName: '省市',
        syntax: ['province_city', '随机生成一个省市名称'],
        example: ['province_city', '-> 陕西省西安市'],
      },
      {
        searchKey: '',
        name: 'province_city_district',
        cnName: '省市区',
        syntax: ['province_city_district', '随机生成一个省市区名称'],
        example: ['province_city_district', '-> 陕西省西安市新城区'],
      },
      {
        searchKey: '',
        name: 'zipcode',
        cnName: '邮政编码',
        syntax: ['zipcode', '随机生成一个邮政编码'],
        example: ['zipcode', '-> 710000'],
      },
      {
        searchKey: '',
        name: 'date',
        cnName: '日期',
        allow: { regexp: createDateRegExp('date') },
        syntax: ['date', '随机生成一个(Y-M-D)格式的日期', 'date|format', '随机生成一个指定格式的日期，format为(YyMmDd)的排列组合'],
        example: ['date', '-> 2020-03-17', 'date|Y-M-D', '-> 2020-03-17', 'date|y-m-d', '-> 20-3-17', 'date|Y/M/D', '-> 2020/03/17', 'date|y/m/d ', '-> 20/3/17'],
      },
      {
        searchKey: '',
        name: 'time',
        cnName: '时间',
        allow: { regexp: createTimeRegExp('time') },
        syntax: ['time', '随机生成一个(H:I:S)格式的时间', 'time|format', '随机生成一个指定格式的时间，format为(HhIiSs)的排列组合'],
        example: ['time', '-> 09:00:00', 'time|H:I:S', '-> 09:00:00', 'time|h:i:s', '-> 9:0:0'],
      },
      {
        searchKey: '',
        name: 'datetime',
        cnName: '日期时间',
        allow: { regexp: createDateTimeRegExp('datetime') },
        syntax: ['datetime', '随机生成一个(Y-M-D H:I:S)格式的日期时间', 'datetime|format', '随机生成一个指定格式的日期时间，format为(YyMmDd HhIiSs)的排列组合'],
        example: [
          'datetime',
          '-> 2020-03-17 09:00:00',
          'datetime|Y-M-D H:I:S',
          '-> 2020-03-17 09:00:00',
          'datetime|y-m-d h:i:s',
          '-> 20-3-17 9:0:0',
          'datetime|Y/M/D H:I:S',
          '-> 2020/03/17 09:00:00',
          'datetime|y/m/d h:i:s',
          '-> 20/3/17 9:0:0',
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
        name: 'dataimage',
        cnName: '图片数据',
        allow: {
          isSwap: true,
          range: { min: 20, max: 1024, minActionText: '图片宽度', maxActionText: '图片高度' },
          parse: parseImage,
          oneOfTypes: ['jpeg', 'png'],
          regexp: createImageRegExp('dataimage'),
        },
        syntax: [
          'dataimage',
          '生成一个200*150的jpeg Base64编码的图片',
          'dataimage|width*height',
          '生成一个指定宽高的jpeg Base64编码的图片',
          'dataimage|width*height,type',
          '生成一个指定宽高和类型的Base64编码的图片，类型type只支持jpeg和png',
        ],
        example: [
          'dataimage',
          '-> data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAM...',
          'dataimage|200*200',
          '-> data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAM...',
          'dataimage|600*400,png',
          '-> data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAM...',
        ],
      },
      {
        searchKey: 'image',
        name: 'imageurl',
        cnName: '图片链接',
        allow: {
          isSwap: true,
          range: { min: 20, max: 1024, minActionText: '图片宽度', maxActionText: '图片高度' },
          parse: parseImage,
          oneOfTypes: ['jpeg', 'png'],
          regexp: createImageRegExp('imageurl'),
        },
        syntax: [
          'imageurl',
          '生成一个200*150的jpeg图片链接',
          'imageurl|width*height',
          '生成一个指定宽高的jpeg图片链接',
          'imageurl|width*height,type',
          '生成一个指定宽高和类型的图片链接，类型type只支持jpeg和png',
        ],
        example: [
          'imageurl',
          '-> http://mock.apicat.net/image/200x150.jpeg',
          'imageurl|200*200',
          '-> http://mock.apicat.net/image/200x200.jpeg',
          'imageurl|200*200,png',
          '-> http://mock.apicat.net/image/200x200.png',
        ],
      },
      {
        searchKey: 'file',
        name: 'fileurl',
        cnName: '文件链接',
        allow: { regexp: createOneOfRegExp('fileurl', ['word', 'excel', 'csv', 'md']) },
        syntax: ['fileurl', '生成一个markdown的文件链接', 'fileurl|type', '生成一个指定类型的文件链接，类型type只支持word、excel、csv、md'],
        example: ['fileurl', '-> http://mock.apicat.net/file/welcome_to_use_apicat.md', 'fileurl|csv', '-> http://mock.apicat.net/file/welcome_to_use_apicat.csv'],
      },
    ],
    int: [
      {
        searchKey: '',
        name: 'int',
        cnName: '整数',
        allow: { range: { min: -1000000, max: 1000000 }, regexp: createNumberRangeRegExp('int') },
        syntax: ['int', '随机生成0-1000的整数', 'int|count', '生成指定数值的整数', 'int|min~max', '随机生成大于等于min小于等于max的整数'],
        example: ['int', '-> 123', 'int|123', '-> 123', 'int|1~10', '-> 5', 'int|-100~-50', '-> -66'],
      },
      {
        searchKey: '',
        name: 'mobile',
        cnName: '手机号',
        syntax: ['mobile', '随机生成一个国内手机号码'],
        example: ['mobile', '-> 13333333333'],
      },
      {
        searchKey: '',
        name: 'idcard',
        cnName: '身份证号',
        syntax: ['idcard', '随机生成一个身份证号码'],
        example: ['idcard', '-> 610102202003170019'],
      },
      {
        searchKey: '',
        name: 'zipcode',
        cnName: '邮政编码',
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
          range: { min: 0, max: 100 },
          parse: parseBoolean,
          regexp: createOneOfRegExp('boolean', ['true', 'false', RE_STR_NORMAL_NUMBER]),
        },
        syntax: [
          'boolean',
          '随机生成一个50%概率的true or false',
          'boolean|value',
          'value为true or false，代表生成指定的布尔值',
          'boolean|probability',
          'probability为0-100的整数，代表生成true的概率，1表示1%，99表示99%',
        ],
        example: [
          'boolean',
          '-> true',
          'boolean|true',
          '-> true',
          'boolean|false',
          '-> false',
          'boolean|0',
          '-> false',
          'boolean|25',
          '-> false',
          'boolean|99',
          '-> true',
          'boolean|100',
          '-> true',
        ],
      },
    ],
    array: [
      {
        searchKey: '',
        name: 'array',
        cnName: '数组',
        allow: { range: { min: 0, max: 50, minActionText: '最小长度', maxActionText: '最大长度' }, regexp: createRangeRegExp('array') },
        syntax: ['array', '数组内的所有元素，随机循环1-5次', 'array|count', '数组内的所有元素循环count次', 'array|min-max', '数组内的所有元素，随机循环大于等于min小于等于max次'],
        example: [
          '// 数组arr1',
          'arr1 = [string|1]',
          '// 数组arr1的规则及结果',
          'array',
          '-> ["a", "b", "c"]',
          '',
          '// 数组arr2',
          'arr2 = [string|1, number|1-9]',
          '// 数组arr2的规则及结果',
          'array|3',
          '-> ["a", 1, "z", 8, "y", 6]',
          '',
          '// 数组arr3',
          'arr3 = [string|1, number|1-9]',
          '// 数组arr3的规则及结果',
          'array|1-3',
          '-> ["a", 1, "z", 8]',
        ],
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
    arrayobject: [
      {
        searchKey: '',
        name: 'array_object',
        alias: 'arrayObject',
        cnName: '数组对象',
        allow: { range: { min: 0, max: 50, minActionText: '最小长度', maxActionText: '最大长度' }, regexp: createRangeRegExp('array_object') },
        syntax: [
          'array_object',
          '数组内的所有元素以key:val的方式组成object，随机生成大于等于1小于等于5个对象元素',
          'array_object|count',
          '数组内的所有元素以key:val的方式组成object，生成count个对象元素',
          'array_object|min-max',
          '数组内的所有元素以key:val的方式组成object，随机生成大于等于min小于等于max个对象元素',
        ],
        example: [
          '// 数组arr1',
          'arr1 = ["username": string|6, "password": string|6]',
          '// 数组arr1的规则及结果',
          'array_object',
          '-> [{"username": "abcdef", "password": "123456"}, {"username": "qwerty", "password": "1q2w3e"}]',
          '',
          '// 数组arr2',
          'arr2 = ["username": string|6]',
          '// 数组arr2的规则及结果',
          'array_object|3',
          '-> [{"username": "abcdef"}, {"username": "qwerty"}, {"username": "ijnuhb"}]',
          '',
          '// 数组arr3',
          'arr3 = ["username": string|6, "password": string|6]',
          '// 数组arr3的规则及结果',
          'array_object|1-3',
          '-> [{"username": "abcdef", "password": "123456"}, {"username": "qwerty", "password": "1q2w3e"}, {"username": "ijnuhb", "password": "5c6v7b"}]',
        ],
      },
    ],
    file: [
      {
        searchKey: '',
        name: 'file',
        cnName: '文件',
        allow: { regexp: createOneOfRegExp('file', ['word', 'excel', 'csv', 'md']) },
        syntax: ['file', '生成一个Markdown文件，以文件流的方式返回', 'file|type', '生成一个指定类型的文件，以文件流的方式返回，类型type只支持word、excel、csv、md'],
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
          oneOfTypes: ['jpeg', 'png'],
          regexp: createImageRegExp('image'),
        },
        syntax: [
          'image',
          '生成一个200*150的jpeg图片流',
          'image|width*height',
          '随机生成一个指定宽高的jpeg图片流',
          'image|width*height,type',
          '随机生成一个指定宽高和类型的图片流，类型type只支持jpeg和png',
        ],
        example: ['image', '-> 返回一个200*150的jpeg图片流', 'image|200*200', '-> 返回一个200*200的jpeg图片流', 'image|600*400,png', '-> 返回一个600*400的png图片流'],
      },
    ],
    float: [
      {
        searchKey: '',
        name: 'float',
        cnName: '浮点数',
        allow: {
          prefix: '整数',
          actionText: '整数',
          range: { min: -1000000, max: 1000000 },
          decimal: { min: 1, max: 10 },
          regexp: creatFloatRangeRegExp('float'),
        },
        syntax: [
          'float',
          '随机生成0-1000的整数和1-3位随机小数的浮点数',
          'float|count.dcount',
          '生成指定数值的整数，小数部分保留dcount位的随机小数',
          'float|count.dmin-dmax',
          '生成指定数值的整数，小数部分保留大于等于dmin小于等于dmax位的随机小数',
          'float|min~max.dmin-dmax',
          '随机生成大于等于min小于等于max的整数，小数部分保留大于等于dmin小于等于dmax位的随机小数',
        ],
        example: ['float', '-> 123.456', 'float|3.2', '-> 3.14', 'float|123.1-3', '-> 123.4', 'float|1~10.1-2', '-> 3.1'],
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
