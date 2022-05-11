<?php

namespace App\Modules\Mock\Parser\Maker;

/**
 * 生成mock数据
 */
class MockRouter
{
    /**
     * 根据规则生成数据
     *
     * @param array $rule mock规则
     * @return int|string|array|object
     */
    public static function generateData($rule)
    {
        switch ($rule['mock_type']) {
            case 'string':
                return StringMock::generate($rule['mock_rule']);
            case 'int':
                return IntMock::generate($rule['mock_rule']);
            case 'float':
                $number = FloatMock::generate($rule['mock_rule']);
                return (float)$number;
            case 'boolean':
                return BooleanMock::generate($rule['mock_rule']);
            case 'array':
                return ArrayMock::generate($rule['mock_rule'], $rule['sub_params']);
            case 'array_object':
                return ArrayObjectMock::generate($rule['mock_rule'], $rule['sub_params']);
            case 'object':
                return ObjectMock::generate($rule['sub_params']);
            case 'paragraph':
                return ParagraphMock::generate($rule['mock_rule']);
            case 'sentence':
                return SentenceMock::generate($rule['mock_rule']);
            case 'word':
                return WordMock::generate($rule['mock_rule']);
            case 'title':
                return TitleMock::generate($rule['mock_rule']);
            case 'cparagraph':
                return CparagraphMock::generate($rule['mock_rule']);
            case 'csentence':
                return CsentenceMock::generate($rule['mock_rule']);
            case 'cword':
                return CwordMock::generate($rule['mock_rule']);
            case 'ctitle':
                return CtitleMock::generate($rule['mock_rule']);
            case 'mobile':
                $mobile = MobileMock::generate();
                if ($rule['type'] == 'int') {
                    $mobile = (int)$mobile;
                }

                return $mobile;
            case 'phone':
                return PhoneMock::generate();
            case 'idcard':
                $idcard = IdcardMock::generate();
                if ($rule['type'] == 'int') {
                    $idcard = (int)$idcard;
                }

                return $idcard;
            case 'url':
                return UrlMock::generate();
            case 'domain':
                return DomainMock::generate();
            case 'ip':
                return IpMock::generate();
            case 'email':
                return EmailMock::generate();
            case 'province':
                return RegionMock::province();
            case 'city':
                return RegionMock::city();
            case 'province_city':
                return RegionMock::provinceCity();
            case 'province_city_district':
                return RegionMock::provinceCityDistrict();
            case 'zipcode':
                $zipcode = RegionMock::zipcode();
                if ($rule['type'] == 'int') {
                    $zipcode = (int)$zipcode;
                }

                return $zipcode;
            case 'date':
                return DateMock::generate($rule['mock_rule']);
            case 'time':
                return TimeMock::generate($rule['mock_rule']);
            case 'datetime':
                return DatetimeMock::generate($rule['mock_rule']);
            case 'timestamp':
                $timestamp = TimestampMock::generate();
                if ($rule['type'] == 'string') {
                    $timestamp = (string)$timestamp;
                }

                return $timestamp;
            case 'image':
                return ImageMock::stream($rule['mock_rule']);
            case 'dataimage':
                return ImageMock::base64($rule['mock_rule']);
            case 'imageurl':
                return ImageMock::url($rule['mock_rule']);
            case 'file':
                return FileMock::stream($rule['mock_rule']);
            case 'fileurl':
                return FileMock::url($rule['mock_rule']);
        }
    }
}
