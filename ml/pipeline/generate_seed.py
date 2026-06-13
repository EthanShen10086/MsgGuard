#!/usr/bin/env python3
"""Generate seed datasets for MsgGuard ML pipeline."""
import csv
import json
from pathlib import Path

ROOT = Path(__file__).parent.parent
SEED = ROOT / "datasets" / "seed"
BENCH = ROOT / "datasets" / "benchmark"

ZH_SPAM = [
    "恭喜您中奖100万请点击链接领取", "免费贷款无抵押当天放款", "您的账户异常请立即验证",
    "兼职刷单日赚500加微信", "澳门赌场开户送888", "低价出售发票可开各类",
    "代开发票联系电话", "您的快递丢失请补填信息", "领导急需用钱请转账",
    "积分即将过期点击兑换", "裸聊加QQ号码", "办理各类证件无需考试",
    "股票内幕消息稳赚", "网贷额度已批请确认", "您的医保卡停用请更新",
    "法院传票请查看", "恭喜成为幸运用户", "零首付购车名额有限",
    "刷单返利先垫付", "加群领红包", "您的社保账户冻结", "高薪招聘日结",
    "代孕服务咨询", "赌博平台注册送彩金", "您的信用卡逾期",
    "免费领取iPhone", "点击链接领现金", "您的包裹被扣请缴费",
    "贷款利息低至1厘", "加微信看片", "您的会员即将到期续费",
    "投资理财年化30%", "您的账号涉嫌违规", "低价代练游戏",
    "办证刻章电话", "您的ETC异常", "网络赌博平台",
    "您的孩子被选为幸运星", "免费领流量", "您的账户存在风险",
    "代写论文包过", "您的积分可兑换现金", "加群领福利",
    "您的医保报销失败", "免费领口罩", "您的贷款已批",
    "点击领取国家补贴", "您的快递无法派送", "免费领鸡蛋",
    "您的账户被冻结解冻", "加微信领红包", "您的会员特权",
    "免费领手机", "您的订单异常", "加QQ领福利",
    "您的账户升级", "免费领话费", "您的包裹滞留",
    "加微信兼职", "您的积分清零", "免费领礼品",
    "您的账户注销", "加群领现金", "您的会员到期",
    "免费领流量包", "您的账户异常登录", "加微信刷单",
    "您的贷款额度", "免费领红包", "您的快递丢失",
    "加QQ兼职", "您的账户风险", "免费领手机壳",
    "您的订单取消", "加微信返利", "您的账户冻结",
    "免费领优惠券", "您的包裹退回", "加群刷单",
    "您的会员续费", "免费领数据线", "您的账户验证",
    "加微信日赚", "您的贷款审批", "免费领耳机",
    "您的快递滞留", "加QQ返利", "您的账户升级通知",
    "免费领充电宝", "您的订单待支付", "加微信领现金",
    "您的积分兑换", "免费领U盘", "您的账户安全",
    "加群日赚", "您的贷款到账", "免费领键盘",
    "您的快递派送", "加微信福利", "您的会员激活",
    "免费领鼠标", "您的账户提醒", "加QQ福利",
    "您的贷款还款", "免费领支架", "您的包裹签收",
    "加微信红包", "您的账户通知", "免费领贴膜",
    "您的订单发货", "加群福利", "您的会员权益",
    "免费领壳", "您的账户变动", "加微信赚钱",
    "您的贷款申请", "免费领线", "您的快递到达",
    "加QQ赚钱", "您的账户余额", "免费领包",
]

ZH_HAM = [
    "您的验证码是123456", "快递已到达菜鸟驿站取件码8888", "妈妈我到了学校",
    "明天下午3点开会请准时", "您的订单已发货顺丰单号SF123", "今晚回家吃饭吗",
    "您的银行卡消费128元", "地铁2号线因故障延误", "您的预约已成功",
    "天气明天有雨记得带伞", "您的外卖已送达", "会议改到402室",
    "您的挂号成功明天9点", "水电费账单已出", "您的航班CA1234已值机",
    "孩子放学了我去接", "您的快递正在派送", "银行余额变动通知",
    "您的预约已确认", "明天体检空腹", "您的订单已签收",
    "地铁1号线正常运行", "您的挂号已取消", "物业费已缴纳",
    "您的航班延误1小时", "周末去爬山吗", "您的外卖正在配送",
    "会议取消通知", "您的预约已改期", "燃气费已扣款",
    "您的火车票已出", "晚上加班不回家", "您的快递已揽收",
    "公交改线通知", "您的挂号提醒", "网费已续费",
    "您的航班已登机", "明天带伞", "您的外卖已取餐",
    "会议记录已发", "您的预约提醒", "电费已缴纳",
    "您的火车票改签", "周末见", "您的快递已入库",
    "地铁恢复运营", "您的挂号成功", "话费已充值",
    "您的航班到达", "记得吃药", "您的外卖已制作",
    "会议纪要", "您的预约成功", "水费已缴纳",
    "您的火车票退票", "生日快乐", "您的快递已出库",
    "公交正常运行", "您的挂号取消", "宽带已续费",
    "您的航班起飞", "注意安全", "您的外卖配送中",
    "工作安排", "您的预约确认", "燃气已缴费",
    "您的火车票成功", "新年快乐", "您的快递派送中",
    "地铁正常", "您的挂号提醒", "流量已充值",
    "您的航班准点", "早点休息", "您的外卖已接单",
    "项目进度", "您的预约改期", "物业已缴费",
    "您的火车票改签成功", "节日快乐", "您的快递已签收",
    "公交正常", "您的挂号成功", "会员已续费",
    "您的航班值机", "多喝水", "您的外卖已送达门口",
    "周报已提交", "您的预约取消", "电费已扣款",
    "您的火车票出票", "周末愉快", "您的快递待取",
    "地铁运行正常", "您的挂号已约", "话费已到账",
    "您的航班登机", "注意保暖", "您的外卖骑手已取餐",
    "任务完成", "您的预约已约", "水费已扣款",
    "您的火车票已取", "工作顺利", "您的快递已放置",
    "公交运营正常", "您的挂号已排", "流量已到账",
    "您的航班落地", "晚安", "您的外卖已放置",
    "日报已发", "您的预约已排", "燃气已扣款",
    "您的火车票已检", "早安", "您的快递已通知",
    "地铁准点", "您的挂号已确", "宽带已到账",
    "您的航班延误", "加油", "您的外卖已通知",
    "月报已交", "您的预约已确", "物业已扣款",
    "您的火车票已登", "辛苦了", "您的快递已提醒",
    "公交准点", "您的挂号已提", "会员已到账",
    "您的航班取消", "谢谢", "您的外卖已提醒",
    "计划已发", "您的预约已提", "电费已到账",
    "您的火车票已退", "不客气", "您的快递已更新",
    "地铁延误恢复", "您的挂号已更", "话费已扣款",
    "您的航班改期", "收到", "您的外卖已更新",
    "总结已写", "您的预约已更", "水费已到账",
    "您的火车票已改", "好的", "您的快递已更新状态",
]

EN_ROWS = [
    ("Free entry in 2 a wkly comp to win FA Cup final tkts 21st May 2005", "spam"),
    ("Go until jurong point, crazy.. Available only in bugis n great world la e buffet", "ham"),
    ("WINNER!! As a valued network customer you have been selected", "spam"),
    ("Had your mobile 11 months or more? U R entitled to Update to the latest", "spam"),
    ("I'm gonna be home soon and i don't want to talk about this stuff anymore tonight", "ham"),
    ("URGENT! You have won a 1 week FREE membership in our £100,000 Prize Jackpot!", "spam"),
    ("I've been searching for the right words to thank you for this breather", "ham"),
    ("SIX chances to win CASH! From 100 to 20,000 pounds txt> CSH11 and send to 87575", "spam"),
    ("Your verification code is 847291", "ham"),
    ("Order #12345 has shipped via UPS tracking 1Z999AA10123456784", "ham"),
    ("Claim your FREE iPhone now click bit.ly/fake", "spam"),
    ("Meeting moved to 3pm in conference room B", "ham"),
    ("Congratulations! You've been selected for a $1000 Walmart gift card", "spam"),
    ("Your appointment is confirmed for Tuesday at 10am", "ham"),
    ("Lowest rate mortgage guaranteed apply now", "spam"),
    ("Dinner at 7? Let me know", "ham"),
    ("You have WON a Nokia 7250i. Call 09061701461 from landline", "spam"),
    ("Your package will arrive tomorrow between 2-4pm", "ham"),
    ("Earn $5000/week working from home", "spam"),
    ("Can you pick up milk on the way home?", "ham"),
    ("URGENT we are trying to contact you. Last weekend's draw shows that you won", "spam"),
    ("Flight AA456 is now boarding at gate 12", "ham"),
    ("Get cheap Cialis online no prescription needed", "spam"),
    ("Thanks for your help today", "ham"),
    ("You are a winner of our £1000 prize draw", "spam"),
    ("Your Uber driver is 2 minutes away", "ham"),
    ("Click here to claim your prize NOW", "spam"),
    ("See you at the gym tomorrow morning", "ham"),
    ("Your account will be suspended verify immediately", "spam"),
    ("The report is ready for review", "ham"),
]

# Extend EN to ~200 rows by duplicating with variations
while len(EN_ROWS) < 200:
    for text, label in list(EN_ROWS)[:30]:
        if len(EN_ROWS) >= 200:
            break
        EN_ROWS.append((text + f" ref{len(EN_ROWS)}", label))


def write_csv(path: Path, rows, fieldnames=("text", "label")):
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=fieldnames)
        w.writeheader()
        for r in rows:
            if isinstance(r, tuple):
                w.writerow({"text": r[0], "label": r[1]})
            else:
                w.writerow(r)


def write_jsonl(path: Path, rows):
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8") as f:
        for r in rows:
            f.write(json.dumps(r, ensure_ascii=False) + "\n")


def main():
    write_csv(SEED / "zh_spam.csv", [(t, "spam") for t in ZH_SPAM])
    write_csv(SEED / "zh_ham.csv", [(t, "ham") for t in ZH_HAM])
    write_csv(SEED / "en_spam_ham.csv", EN_ROWS)

    zh_test = [{"text": t, "label": "spam"} for t in ZH_SPAM[:25]] + [{"text": t, "label": "ham"} for t in ZH_HAM[:25]]
    en_test = [{"text": t, "label": l} for t, l in EN_ROWS[:50]]
    adversarial = [
        {"text": "您的验证码是847291，5分钟内有效", "label": "ham"},
        {"text": "【顺丰快递】您的包裹已到达，取件码1234", "label": "ham"},
        {"text": "Your verification code is 123456", "label": "ham"},
        {"text": "Amazon: Your order has shipped", "label": "ham"},
        {"text": "验证码1234请勿泄露", "label": "ham"},
        {"text": "快递取件码8888菜鸟驿站", "label": "ham"},
        {"text": "【银行】您尾号1234消费128元", "label": "ham"},
        {"text": "Meeting at 3pm room 402", "label": "ham"},
        {"text": "免费验证码领取点击链接", "label": "spam"},
        {"text": "您的验证码中奖请点击", "label": "spam"},
        {"text": "快递丢失请点击补填", "label": "spam"},
        {"text": "Verify account click here free", "label": "spam"},
        {"text": "Your package held pay fee click", "label": "spam"},
        {"text": "【系统】验证码123456", "label": "ham"},
        {"text": "OTP 847291 valid 5 min", "label": "ham"},
        {"text": "积分兑换验证码8888", "label": "ham"},
        {"text": "登录验证码请勿转发", "label": "ham"},
        {"text": "Your code is 999888 do not share", "label": "ham"},
        {"text": "恭喜验证码中奖", "label": "spam"},
        {"text": "Free gift verify now click", "label": "spam"},
    ]
    write_jsonl(BENCH / "test_zh.jsonl", zh_test)
    write_jsonl(BENCH / "test_en.jsonl", en_test)
    write_jsonl(BENCH / "adversarial.jsonl", adversarial)
    print(f"Generated seed: zh_spam={len(ZH_SPAM)}, zh_ham={len(ZH_HAM)}, en={len(EN_ROWS)}")


if __name__ == "__main__":
    main()
